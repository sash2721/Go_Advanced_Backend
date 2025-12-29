package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Time struct {
	CurrentTime string `json:"currentTime"`
	RequestID   string `json:"requestID"`
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieving the requestID
	requestId := r.Context().Value("requestID").(string)

	currentTime := time.Now() // find current time

	var timeFormat string = "2006/01/02, 03:04 PM"
	formattedTime := currentTime.Format(timeFormat) // formatting time in string

	jsonString := Time{CurrentTime: formattedTime, RequestID: requestId} // creating a json string
	jsonData, err := json.Marshal(jsonString)                            // converting in JSON

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
		"Time API Response Sent!",
		slog.String("RequestID", requestId),
		slog.String("CurrentTime", formattedTime),
	)
}
