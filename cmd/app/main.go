package main

import (
	"crud-golang/internal/app"
	"crud-golang/internal/config"
	"log"
)

func main() {
	// Configuration
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Run
	app.Run()
}
