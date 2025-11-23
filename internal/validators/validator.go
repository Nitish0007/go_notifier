package validators

import (
	"fmt"
	"strings"
)

// ValidationError represents a single field validation error
// As it contains the Error method
type ValidationError struct {
	Message string
	Field   string
	Tag     string // The validation tag that failed (e.g., "required", "email")
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) GetField() string {
	return e.Field
}

func (e *ValidationError) GetTag() string {
	return e.Tag
}

// ValidationErrors is a collection of multiple validation errors.
// It implements the error interface and provides methods to access individual errors.
type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		// When there are no validation errors, it's more appropriate to return an empty string or a message indicating no error.
		return ""
	}

	if len(e) == 1 {
		return e[0].Error()
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(messages, "; "))
}

// GetErrors returns the slice of validation errors
func (e ValidationErrors) GetErrors() []*ValidationError {
	return []*ValidationError(e)
}

// GetFieldErrors returns all errors for a specific field
func (e ValidationErrors) GetFieldErrors(fieldName string) []*ValidationError {
	var fieldErrors []*ValidationError
	for _, err := range e {
		if err.Field == fieldName {
			fieldErrors = append(fieldErrors, err)
		}
	}
	return fieldErrors
}

// HasField checks if there are any errors for a specific field
func (e ValidationErrors) HasField(fieldName string) bool {
	for _, err := range e {
		if err.Field == fieldName {
			return true
		}
	}
	return false
}
