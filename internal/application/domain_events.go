package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/database"
	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/internal/logger"
)

type eventService interface {
	Create(ctx context.Context, event events.Event) error
	Get(ctx context.Context, id uuid.UUID) (events.Event, error)
	GetList(ctx context.Context, from time.Time, to time.Time) ([]events.Event, error)
	Update(ctx context.Context, id uuid.UUID, event events.UpdateEvent) error
	Delete(ctx context.Context, id string) error
	GetRecentEvents(ctx context.Context, from time.Time) (*[]events.Event, error)
	DeleteOldEvents(ctx context.Context, date time.Time) (bool, error)
}

type EventsDomain struct {
	eventsService eventService
}

func buildEventsDomain(db database.DBConnection, logger logger.Logger) *EventsDomain {
	repo := events.NewEventsRepository(db)
	eventsService := events.NewEventsService(repo, logger)
	return &EventsDomain{
		eventsService: eventsService,
	}
}
