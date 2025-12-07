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

type Handlers struct {
	eventService eventService
	logger       logger.Logger
}

func NewHandlers(logger logger.Logger, eventService eventService) *Handlers {
	return &Handlers{
		eventService: eventService,
		logger:       logger,
	}
}

func NewRouter(handlers *Handlers) *chi.Mux {
	router := chi.NewMux()

	// Configure middleware
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	
	// Configure routes
	router.Get("/healthcheck", handlers.healthcheck)
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/events", handlers.createEvent)
		r.Get("/events/{id}", handlers.getEvent)
		r.Get("/events", handlers.getEventsList)
		r.Put("/events/{id}", handlers.updateEvent)
		r.Delete("/events", handlers.deleteEvent)
	})

	return router
}

