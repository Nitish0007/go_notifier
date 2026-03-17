package api

import (
	"net/http"
	"errors"
	"encoding/json"
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

