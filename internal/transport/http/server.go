package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/logger"
)

type Server struct {
	server   *http.Server
	listener net.Listener
	logger   logger.Logger
}

const ServerShutdownTimeout = 10 * time.Second

// create a new HTTP server instance
// TODO: add - handler http.Handler
func NewServer(cfg *config.HTTPConfig, logger logger.Logger) (*Server, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	server := &http.Server{}

	return &Server{
		server:   server,
		listener: listener,
		logger:   logger,
	}, nil
}

// start serving HTTP requests
func (s *Server) Start(cfg *config.HTTPConfig) error {
	s.logger.Infof("Server started on port :%v", cfg.Port)

	// serving HTTP requests should be in a separate goroutine cause 
	// Serve() blocks the main thread
	go func() {
		err := s.server.Serve(s.listener)
		// Panic and stop the entire application in case of a server start error.
		// This is a rare case, as we've already started listening on the port.
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Errorf("Error during server shutdown: %v", err)
		return err
	}

	s.logger.Info("Server successfully shutted down")
	return nil
}
