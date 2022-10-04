package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}
}

func main() {
	if err := NewCli().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
