package services

import (
	"context"
	"io"
	"log/slog"
	"net/http"
)

type EchoService struct{}

func NewEchoService() *EchoService {
	return &EchoService{}
}

func (s *EchoService) EchoResponse(r *http.Request) ([]byte, string, context.Context, error) {
	// extracting the requestId from the request context
	requestId := r.Context().Value("requestID").(string)

	// Creating the request context
	ctx := r.Context()

	// read the request body
	response, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Error(
			"Error while reading the Request Body",
			slog.Any("Error", err),
			slog.String("RequestID", requestId),
		)
		return nil, requestId, nil, err
	}
	defer r.Body.Close() // close the connection at the end

	return response, requestId, ctx, err
}
