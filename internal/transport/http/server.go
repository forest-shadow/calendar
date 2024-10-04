package http

import (
	"github.com/forest-shadow/calendar/internal/config"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	server *http.Server
	listener net.Listener
}

// create a new HTTP server instance
// TODO: add - handler http.Handler
func NewServer(cfg *config.HTTPConfig) (*Server, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	server := &http.Server{}

	return &Server{
		server: server,
		listener: listener,
	}, nil
}

// start serving HTTP requests
func (s *Server) Start(cfg *config.HTTPConfig) error {
	fmt.Printf("Server started on port :%v\n", cfg.Port)
	err := s.server.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}
