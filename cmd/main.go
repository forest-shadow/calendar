package main

import (
	"fmt"
	"log"
	"github.com/forest-shadow/calendar/internal/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}
	fmt.Printf("config: %+v\n", cfg)
}
