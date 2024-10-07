package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/forest-shadow/calendar/internal/application"
)

func main() {
	// Create a context that will be canceled when the 
	// program receives an interrupt or termination signal
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
