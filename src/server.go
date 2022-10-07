package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
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

func StartWebServer() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	app := gin.New()

	if os.Getenv("LOGGER") != "false" {
		app.Use(gin.Logger())
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: MustGetEnv("SENTRY_DSN", ""),
	})
	if err != nil {
		log.Println("Sentry initialization failed: " + err.Error())
	} else {
		app.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

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

	RegisterRoutes(app)

	srv := &http.Server{
		Addr:    MustGetEnv("HOST", "0.0.0.0") + ":" + MustGetEnv("PORT", "3333"),
		Handler: app,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)

	<-sigs

	log.Println("Closing web server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
