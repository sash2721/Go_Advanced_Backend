package services

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type TimeService struct{}

type Time struct {
	CurrentTime string `json:"currentTime"`
	RequestID   string `json:"requestID"`
}

func NewTimeService() *TimeService {
	return &TimeService{}
}

func (*TimeService) GetTime(r *http.Request) ([]byte, string, string, error) {
	// Retrieving the requestID
	requestId := r.Context().Value("requestID").(string)

	currentTime := time.Now() // find current time

	var timeFormat string = "2006/01/02, 03:04 PM"
	formattedTime := currentTime.Format(timeFormat) // formatting time in string

	jsonString := Time{CurrentTime: formattedTime, RequestID: requestId} // creating a json string
	jsonData, err := json.Marshal(jsonString)                            // converting in JSON

	if err != nil {
		slog.Error(
			"Error while Marshaling the Time API response",
			slog.Any("Error", err),
			slog.String("RequestID", requestId),
		)
		return nil, "", "", err
	}

	return jsonData, formattedTime, requestId, nil
}
