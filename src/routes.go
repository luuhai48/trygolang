package main

import (
	"net/http"

	_ "luuhai48/trygolang/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(app *gin.Engine) {
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Ok",
		})
	})

	app.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
