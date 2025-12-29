package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Health struct {
	Message   string `json:"message"`
	RequestID string `json:"requestID"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieving the requestID to add in the Response
	requestId := r.Context().Value("requestID").(string)

	jsonString := Health{
		Message:   "healthy",
		RequestID: requestId,
	}
	jsonData, err := json.Marshal(jsonString)

	if err != nil {
		slog.Error(
			"Error while Marshaling the Health API response",
			slog.Any("Error", err),
			slog.String("RequestID", requestId),
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error while processing the data!"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	slog.Info(
		"Health API Response Sent!",
		slog.String("RequestID", requestId),
	)
}
