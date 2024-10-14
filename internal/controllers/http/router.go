package handlers

import (
	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"

	"github.com/forest-shadow/calendar/internal/logger"
)

type handlers struct {
	logger logger.Logger
}

func NewRouter(
	logger logger.Logger,
) *chi.Mux {
	router := chi.NewMux()

	handlers := handlers{
		logger: logger,
	}
	handlers.build(router)

	return router
}

func (h *handlers) build(router chi.Router) {
	router.Use(middleware.Recoverer)
	router.Get("/healthcheck", h.healthcheck)
}
