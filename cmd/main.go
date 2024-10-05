package main

import (
	"log"

	"github.com/forest-shadow/calendar/internal/application"
)

func main() {
	if err := application.Run(); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
