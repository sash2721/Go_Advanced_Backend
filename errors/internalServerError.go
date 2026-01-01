package errors

import (
	"encoding/json"
	"net/http"
)

type InternalServerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func NewInternalServerError(message string, err error) ([]byte, *InternalServerError) {
	customError := &InternalServerError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   err,
	}

	jsonData, err := json.Marshal(customError)
	if err != nil {
		return []byte(`{"code":500,"message":"Internal Server Error","error":"Error while marshaling the internal server error"}`), nil
	}

	return jsonData, customError
}
