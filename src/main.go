package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}
}

// @title 			Go API
// @version 		1.0
// @description Go API documentation

// @contact.name   luuhai48
// @contact.url    https://luuhai48.github.io
// @contact.email  luuhai.hn48@gmail.com

// @BasePath 	/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\r- Ctrl+C pressed in Terminal")

		Shutdown()

		os.Exit(0)
	}()

	if err := NewCli().Run(os.Args); err != nil {
		log.Fatal(err)
	}
	log.Println("Exiting")
}
