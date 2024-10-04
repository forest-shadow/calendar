package main

import (
	"log"
	"github.com/forest-shadow/calendar/internal/application"
	"github.com/forest-shadow/calendar/internal/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}
	app := application.NewApp(cfg)
	app.HttpServer.Start(&cfg.HTTP)
}
