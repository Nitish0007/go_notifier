package api

import (
	"log"
	"net/http"
	"encoding/json"
)

type ApiResponse struct {
	StatusCode int     `json:"status,omitempty"`
	Message    string  `json:"message,omitempty"`
	Data       any     `json:"data,omitempty"`
}

type ApiErrorResponse struct {
	StatusCode int     `json:"status,omitempty"`
	Error      string  `json:"error,omitempty"`
}

func (r *ApiResponse) WriteResponse(w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("Unable to encode JSON: %v", err)
	}
}

func (r *ApiErrorResponse) WriteErrorResponse(w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("Unable to encode JSON: %v", err)
	}
}