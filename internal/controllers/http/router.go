package handlers

import (
	// "net/http"

	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"

	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/internal/logger"
)

type handlers struct {
	eventService *events.EventService
	logger       logger.Logger
}

func NewRouter(
	logger logger.Logger,
	eventService *events.EventService,
) *chi.Mux {
	router := chi.NewMux()

	handlers := handlers{
		eventService: eventService,
		logger:       logger,
	}
	handlers.build(router)

	return router
}

func (h *handlers) build(router chi.Router) {
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Get("/healthcheck", h.healthcheck)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/events", h.createEvent)
		r.Get("/events/{id}", h.getEvent)
		r.Get("/events", h.getEventsList)
		r.Put("/events/{id}", h.updateEvent)
		r.Delete("/events/{id}", h.deleteEvent)
	})
}
