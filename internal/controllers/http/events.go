package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/apperrs"
	"github.com/forest-shadow/calendar/internal/events"
)

func (h *handlers) createEvent(w http.ResponseWriter, r *http.Request) {
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

func (h *handlers) getEvent(w http.ResponseWriter, r *http.Request) {
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

func (h *handlers) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID := chi.URLParam(r, "id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	var event events.EventUpdateDTO
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
	err = h.eventService.Update(ctx, id, event)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) deleteEvent(w http.ResponseWriter, r *http.Request) {
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
