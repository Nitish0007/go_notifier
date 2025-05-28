package utils

import (
	"log"
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

func WriteJSONResponse(w http.ResponseWriter, status int, data any, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := make(map[string]any)
	if message != "" {
		response["message"] = message
	}

	if data != nil {
		response["data"] = data
	}

	if(len(response) == 0){
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal("Unable to encode JSON")
		return err
	}
	return nil
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		return err
	}
	return nil
}