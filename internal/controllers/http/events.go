package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/apperrs"
	"github.com/forest-shadow/calendar/internal/events"
)

type eventUpdateDTO struct {
	Title        *string    `json:"title,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Description  *string    `json:"description,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	NotifyBefore *string    `json:"notify_before,omitempty"`
}

func (e *eventUpdateDTO) toDomain() events.UpdateEvent {
	result := make(events.UpdateEvent)
	if e.Title != nil {
		result["title"] = *e.Title
	}
	if e.StartTime != nil {
		result["start_time"] = *e.StartTime
	}
	if e.EndTime != nil {
		result["end_time"] = *e.EndTime
	}
	if e.Description != nil {
		result["description"] = *e.Description
	}
	if e.UserID != nil {
		result["user_id"] = *e.UserID
	}
	if e.NotifyBefore != nil {
		result["notify_before"] = *e.NotifyBefore
	}
	return result
}

type eventListFilterDTO struct {
	From time.Time `query:"from"`
	To   time.Time `query:"to"`
}

func (h *Handlers) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newEvent events.Event

	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	err = r.Body.Close()
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	if newEvent.NotifyBefore != nil {
		if err := validateISO8601Duration(*newEvent.NotifyBefore); err != nil {
			h.handleError(ctx, w, fmt.Errorf("%w: %v", apperrs.ErrConditionViolation, err))
			return
		}
	}

	err = h.eventService.Create(ctx, newEvent)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) getEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID := chi.URLParam(r, "id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	event, err := h.eventService.Get(ctx, id)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	_, err = w.Write(eventJSON)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) getEventsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var eventsDto eventListFilterDTO
	err := json.NewDecoder(r.Body).Decode(&eventsDto)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	err = r.Body.Close()
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	events, err := h.eventService.GetList(ctx, eventsDto.From, eventsDto.To)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	eventsJSON, err := json.Marshal(events)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	_, err = w.Write(eventsJSON)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID := chi.URLParam(r, "id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	var event eventUpdateDTO
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	err = r.Body.Close()
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	if event.NotifyBefore != nil {
		if err := validateISO8601Duration(*event.NotifyBefore); err != nil {
			h.handleError(ctx, w, fmt.Errorf("%w: %v", apperrs.ErrConditionViolation, err))
			return
		}
	}
	err = h.eventService.Update(ctx, id, event.toDomain())
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID := chi.URLParam(r, "id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	err = h.eventService.Delete(ctx, id.String())
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
