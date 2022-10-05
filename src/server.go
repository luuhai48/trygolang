package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

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

	app.Use(
		gin.CustomRecovery(customRecoveryHandler),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	if os.Getenv("LOGGER") != "false" {
		app.Use(gin.Logger())
	}

	RegisterRoutes(app)

	return app
}
