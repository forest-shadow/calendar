package application

import (
	"github.com/forest-shadow/calendar/internal/database"
	"github.com/forest-shadow/calendar/internal/events"
)

type EventsDomain struct {
	eventsService *events.EventService
}

func buildEventsDomain(db database.DBConnection) *EventsDomain {
	repo := events.NewEventsRepository(db)
	eventsService := events.NewEventService(repo)
	return &EventsDomain{
		eventsService: eventsService,
	}
}
