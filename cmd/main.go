package main

import (
	"log"
	"solution/cmd/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	application.Run()
}
