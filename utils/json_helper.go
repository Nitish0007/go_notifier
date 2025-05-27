package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseJSONBody(r *http.Request) (map[string]any, error) {
	var payload map[string]any
	if r.Body == nil {
		return nil, errors.New("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if len(payload) == 0 {
		return nil, errors.New("request body is empty or invalid")
	}
	return payload, nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}