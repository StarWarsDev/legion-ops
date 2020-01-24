package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	localFilePath := os.Getenv("CLIENT_FILES_PATH")
	if localFilePath == "" {
		localFilePath = "./client/build"
	}

	StartServer(port, localFilePath, wait)
}
