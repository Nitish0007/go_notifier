package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func GetPathParam(r *http.Request, paramName string) (string, error){
	pathParams := chi.URLParam(r, paramName)
	if pathParams == "" {
		return "", fmt.Errorf("%s parameter not found in path", paramName)
	}
	return pathParams, nil
}

func GetQueryParams(r *http.Request) (map[string]string, error) {
	queryParams := r.URL.Query()
	queryMap := make(map[string]string)
	for key, values := range queryParams {
		queryMap[key] = values[0]
	}
	return queryMap, nil
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
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