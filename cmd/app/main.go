package main

import (
	"github.com/Vaixle/crud-golang/internal/app"
	"github.com/Vaixle/crud-golang/internal/config"
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
