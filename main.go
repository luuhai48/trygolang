package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.GET("/api", func(c *gin.Context) {
		firstName := c.Query("firstName")
		lastName := c.Query("lastName")
		c.String(http.StatusOK, "Hi %s %s", firstName, lastName)
	})

	r.Run(":3333")
}
