package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// check if the request is a POST request or not
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Only POST methods allowed"}`))
		return
	}

	// extracting the requestId from the request context
	requestId := r.Context().Value("requestID").(string)

	// read the request body
	body, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Error(
			"Error while reading the Request Body",
			slog.Any("Error", err),
			slog.String("RequestID", requestId),
		)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Failed to read the request body!"}`))
		return
	}
	defer r.Body.Close() // close the connection at the end

	// read the context
	ctx := r.Context()

	select {
	case <-time.After(2 * time.Second):
		// return the echo json back
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)

		slog.Info(
			"Echo API Response Sent!",
			slog.String("RequestID", requestId),
			slog.String("Body", string(body)),
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
