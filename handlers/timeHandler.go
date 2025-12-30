package handlers

import (
	"advancedBackend/services"
	"log/slog"
	"net/http"
)

type TimeHandler struct {
	Service *services.TimeService
}

func (h *TimeHandler) HandleTimeFunction(w http.ResponseWriter, r *http.Request) {
	// sending the request to service
	responseData, currentTime, requestId, err := h.Service.GetTime(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error while processing the request!"}`))

		slog.Error(
			"Time API Failed!",
			slog.Any("Error:", err),
		)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseData))

		slog.Info(
			"Time API Response Sent!",
			slog.String("RequestID", requestId),
			slog.String("CurrentTime", currentTime),
		)
		return
	}
}
