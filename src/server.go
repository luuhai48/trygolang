package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func customRecoveryHandler(c *gin.Context, err any) {
	message := "Something went wrong"
	if gin.IsDebugging() {
		switch e := err.(type) {
		case string:
			message += ": " + e
		case runtime.Error:
			message += ": " + e.Error()
		case error:
			message += ": " + e.Error()
		}
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func NewServer() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))

	app := gin.New()

	app.Use(gin.CustomRecovery(customRecoveryHandler), cors.Default())

	if os.Getenv("LOGGER") != "false" {
		app.Use(gin.Logger())
	}

	return app
}
