package api

import (
	"net/http"
	"errors"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/Nitish0007/go_notifier/internal/shared/validators"
)

func ParseJSONBody[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, errors.New("request body is empty")
	}

	var req *T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func ParseAndValidateRequest[T any](r *http.Request) (*T, error) {
	req, err := ParseJSONBody[T](r)
	if err != nil {
		return nil, err
	}

	return ValidateRequestData(req)
}

func ValidateRequestData[T any](reqData *T) (*T, error) {
	validator := validators.NewModelValidator[T]()
	if err := validator.ValidateStruct(reqData); err != nil {
		return nil, err
	}
	return reqData, nil
}

func WriteResponse(w http.ResponseWriter, status int, data any, message string) {
	response := ApiResponse{
		StatusCode: status,
		Message: message,
		Data: data,
	}

	WriteResponseHeader(w, status)
	response.WriteResponse(w)
}

func WriteErrorResponse(w http.ResponseWriter, status int, errorMessage string) {
	response := ApiErrorResponse{
		StatusCode: status,
		Error: errorMessage,
	}

	WriteResponseHeader(w, status)
	response.WriteErrorResponse(w)
}

func WriteResponseHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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



