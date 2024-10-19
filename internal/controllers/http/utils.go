package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/forest-shadow/calendar/internal/apperrs"
)

func (h *handlers) handleError(ctx context.Context, w http.ResponseWriter, err error) {
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
