package utils

import (
	"encoding/json"
	"net/http"
)

func JsonErrorUnless(condition bool, message string, details string, status int, w http.ResponseWriter) {
	if !condition {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   message,
			"details": details,
		})
	}
}
