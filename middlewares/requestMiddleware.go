package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUUID, err := exec.Command("uuidgen").Output()

		if err != nil {
			slog.Error(
				"Error while processing the request",
				slog.Any("Error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error while processing the request!"}`))
			return
		}

		// Trim newline or trailing spaces from uuidgen output
		trimmedRequestID := strings.TrimSpace(string(newUUID))

		// Store the UUID in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "requestID", trimmedRequestID)
		slog.Info(
			"Added the RequestID to the context:",
			slog.String("RequestID", trimmedRequestID),
			slog.Any("Context", ctx),
		)

		// Store the requestID in the response Header
		w.Header().Set("X-Request-ID", trimmedRequestID)
		slog.Info(
			"Added the RequestID to the Response Headers:",
			slog.String("RequestID", trimmedRequestID),
		)

		// call the next Handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
