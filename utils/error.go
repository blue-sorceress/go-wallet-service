package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Code             int    `json:"code"`
}

/*
func NewErrorResponse(errorType string, errorDescription string, code int) ErrorResponse {
	return ErrorResponse{
		Error:            errorType,
		ErrorDescription: errorDescription,
		Code:             code,
	}
}*/

func HttpErrorResponse(w http.ResponseWriter, errorType string, errorDescription string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:            errorType,
		ErrorDescription: errorDescription,
		Code:             code,
	})
}
