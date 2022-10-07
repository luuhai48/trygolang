package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

var (
	invalidJSONBodyResponse = gin.H{"message": "Invalid json body"}

	cache   map[string]*gojsonschema.Schema = make(map[string]*gojsonschema.Schema)
	cacheMx sync.Mutex
)

type ErrCannotBuildSchema struct {
	err error
}

func (e *ErrCannotBuildSchema) Error() string {
	return fmt.Sprintf("Cannot build schema: %v", e.err)
}

func NewErrCannotBuildSchema(err error) *ErrCannotBuildSchema {
	return &ErrCannotBuildSchema{err}
}

type ErrSchemaValidation struct {
	Errors []string
}

func (e *ErrSchemaValidation) Error() string {
	return strings.Join(e.Errors, "; ")
}

func NewErrSchemaValidation(errors []string) *ErrSchemaValidation {
	return &ErrSchemaValidation{
		Errors: errors,
	}
}

func buildSchemaFromString(schemaStr string) (*gojsonschema.Schema, error) {
	if sch, found := cache[schemaStr]; !found {
		// if value not found, we should create new schema and put it in cache

		cacheMx.Lock()
		defer cacheMx.Unlock()

		// now read again, probably other goroutine already write value in cache
		if sch, found = cache[schemaStr]; found {
			return sch, nil
		}

		// create new schema
		sch, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaStr))
		if err != nil {
			return nil, NewErrCannotBuildSchema(err)
		}

		cache[schemaStr] = sch
		return sch, nil
	} else {
		return sch, nil
	}
}

func drainHTTPRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, io.EOF
	}

	// read body
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// replace request body with new one
	req.Body = io.NopCloser(bytes.NewBuffer(buf))

	return buf, nil
}

func validateBodyUsingSchema(req *http.Request, schema *gojsonschema.Schema) error {
	body, err := drainHTTPRequestBody(req)
	if err != nil {
		return err
	}

	// validate body
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return err
	}

	// schema not valid, create validation error
	if !result.Valid() {
		var errors []string
		for _, er := range result.Errors() {
			errors = append(errors, er.String())
		}
		return NewErrSchemaValidation(errors)
	}

	return nil
}

func ValidateSchema(handler gin.HandlerFunc, schema *gojsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateBodyUsingSchema(c.Request, schema); err == nil {
			handler(c)
		} else {
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	c.Abort()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		c.JSON(http.StatusBadRequest, invalidJSONBodyResponse)
	} else {
		switch v := err.(type) {
		case *json.SyntaxError:
			c.JSON(http.StatusBadRequest, invalidJSONBodyResponse)
		case *ErrSchemaValidation:
			c.JSON(http.StatusBadRequest, gin.H{
				"messages": v.Errors,
			})
		default:
			c.Status(http.StatusInternalServerError)
		}
	}
}

func Validate(handler gin.HandlerFunc, schemaStr string) gin.HandlerFunc {
	sch, err := buildSchemaFromString(schemaStr)
	if err != nil {
		panic(fmt.Sprintf("Cannot build schema from string %v", schemaStr))
	}
	return ValidateSchema(handler, sch)
}

func BindJSON(c *gin.Context, schemaStr string, obj interface{}) error {
	sch, err := buildSchemaFromString(schemaStr)
	if err != nil {
		panic(err)
	}
	return BindJSONSchema(c, sch, obj)
}

func BindJSONSchema(c *gin.Context, schema *gojsonschema.Schema, obj interface{}) (err error) {
	defer func() {
		if err != nil {
			handleError(c, err)
		}
	}()

	// validate body
	if err = validateBodyUsingSchema(c.Request, schema); err != nil {
		return
	}

	// read body and unmarshal json
	var body []byte
	if body, err = io.ReadAll(c.Request.Body); err != nil {
		return
	}

	if err = json.Unmarshal(body, obj); err != nil {
		return
	}

	return
}
