package handlers

import (
	"net/http"
)

func (h *handlers) healthcheck(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info("healthcheck")
	w.WriteHeader(http.StatusOK)
}
