package errors

import (
	"encoding/json"
	"net/http"
)

type BadRequestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func NewBadRequestError(message string, err error) ([]byte, *BadRequestError) {
	customError := &BadRequestError{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   err,
	}

	jsonData, err := json.Marshal(customError)
	if err != nil {
		return []byte(`{"code":500,"message":"Internal Server Error","error":"Error while marshaling the internal server error"}`), nil
	}

	return jsonData, customError
}
