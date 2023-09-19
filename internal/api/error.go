package api

import (
	"encoding/json"
	"net/http"
)

// //////////////////////////////////////////////////
// encode error

func encodeError(resp http.ResponseWriter, statusCode int, message string) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)

	// try to encode error >>> no need to check error at encoding
	json.NewEncoder(resp).Encode(toJsonError(statusCode, message))
}

func toJsonError(code int, message string) *JsonErrorResponse {
	return &JsonErrorResponse{
		Success: false,
		Error: &JsonError{
			Code:    code,
			Message: message,
		},
	}
}

type JsonErrorResponse struct {
	Success bool       `json:"success,omitempty"`
	Error   *JsonError `json:"error,omitempty"`
}

type JsonError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
