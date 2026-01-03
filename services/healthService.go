package services

import (
	"advancedBackend/errors"
	"encoding/json"
	"log/slog"
	"net/http"
)

type HealthService struct{}

type Health struct {
	Message   string `json:"message"`
	RequestID string `json:"requestID"`
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) GetHealth(r *http.Request) ([]byte, string, error, []byte) {
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
		errorJsonData, internalServerError := errors.NewInternalServerError("Error while Marshaling the Health API response", err)

		return nil, requestId, internalServerError, errorJsonData
	}

	return jsonData, requestId, nil, nil
}
