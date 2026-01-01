package handlers

import (
	"advancedBackend/services"
	"log/slog"
	"net/http"
	"time"
)

type EchoHandler struct {
	Service *services.EchoService
}

func (h *EchoHandler) HandleEchoFunction(w http.ResponseWriter, r *http.Request) {
	// check if the request is a POST request or not
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Only POST methods allowed"}`))
		return
	}

	// calling the service
	responseBody, requestId, ctx, err, errorJsonData := h.Service.EchoResponse(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorJsonData)

		slog.Info(
			"Echo API Failed!",
			slog.String("RequestID", requestId),
			slog.String("Body", string(responseBody)),
		)
		return
	}

	select {
	case <-time.After(2 * time.Second):
		// return the echo json back
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)

		slog.Info(
			"Echo API Response Sent!",
			slog.String("RequestID", requestId),
			slog.String("Body", string(responseBody)),
		)
	case <-ctx.Done():
		slog.Error(
			"Client Disconnected",
			slog.Any("Error", ctx.Err()),
			slog.String("RequestID", requestId),
		)
		return
	}
}
