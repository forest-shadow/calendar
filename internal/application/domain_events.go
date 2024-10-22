package application

import (
	"github.com/forest-shadow/calendar/internal/database"
	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/internal/logger"
)

type EventsDomain struct {
	eventsService events.EventService
}

func buildEventsDomain(db database.DBConnection, logger logger.Logger) *EventsDomain {
	repo := events.NewEventsRepository(db)
	eventsService := events.NewEventService(repo, logger)
	return &EventsDomain{
		eventsService: eventsService,
	}
}
