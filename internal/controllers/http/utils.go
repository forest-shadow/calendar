package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	p "github.com/rickb777/period"

	"github.com/forest-shadow/calendar/internal/apperrs"
)

func (h *Handlers) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	h.logger.With("operation", chi.RouteContext(ctx).RoutePattern()).Error(err.Error())

	switch {
	case errors.Is(err, apperrs.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, apperrs.ErrConditionViolation):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, apperrs.ErrAlreadyExist):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, apperrs.ErrUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		h.logger.Error("write error", err.Error())
		return
	}
}

func validateISO8601Duration(durationStr string) error {
	_, err := p.Parse(durationStr)
	if err != nil {
		return fmt.Errorf("invalid ISO 8601 duration format: %w", err)
	}
	return nil
}
