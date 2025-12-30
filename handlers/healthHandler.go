package handlers

import (
	"advancedBackend/services"
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	Service *services.HealthService
}

func (h *HealthHandler) HandleHealthFunction(w http.ResponseWriter, r *http.Request) {
	// calling the service with the request data
	responseData, requestId, err := h.Service.GetHealth(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error while processing the request"}`))

		slog.Error(
			"Health API Failed!",
			slog.Any("Error", err),
		)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseData))

		slog.Info(
			"Health API Response Sent!",
			slog.String("RequestID", requestId),
		)
	}
}
