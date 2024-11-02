package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"

	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/internal/logger"
)

type eventService interface {
	Create(ctx context.Context, event events.Event) error
	Get(ctx context.Context, id uuid.UUID) (events.Event, error)
	GetList(ctx context.Context, from time.Time, to time.Time) ([]events.Event, error)
	Update(ctx context.Context, id uuid.UUID, event events.UpdateEvent) error
	Delete(ctx context.Context, id string) error
}

type handlers struct {
	eventService eventService
	logger       logger.Logger
}

func NewRouter(
	logger logger.Logger,
	eventService eventService,
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
