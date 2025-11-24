package validators

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ModelValidator is a generic validator that can validate and convert between
// JSON bytes, map[string]any, and struct types for any model.
type ModelValidator[T any] struct {
	validator *validator.Validate
}

// NewModelValidator creates a new instance of ModelValidator for type T.
// The validator instance is reusable and thread-safe.
func NewModelValidator[T any]() *ModelValidator[T] {
	return &ModelValidator[T]{
		validator: validator.New(),
	}
}

// ValidateFromJSON validates raw JSON bytes and returns a validated struct instance.
// This is useful when receiving JSON from HTTP requests.
func (v *ModelValidator[T]) ValidateFromJSON(data []byte) (*T, error) {
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := v.ValidateStruct(&obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

// ValidateFromMap validates a map[string]any and returns a validated struct instance.
// This is useful when working with parsed JSON payloads from handlers.
func (v *ModelValidator[T]) ValidateFromMap(data map[string]any) (*T, error) {
	// Convert map to JSON first, then unmarshal to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	return v.ValidateFromJSON(jsonData)
}

// ValidateStruct validates an existing struct instance.
// This is useful for validating structs that were created programmatically
// or retrieved from the database before use.
func (v *ModelValidator[T]) ValidateStruct(obj *T) error {
	if obj == nil {
		return &ValidationError{
			Message: "struct cannot be nil",
			Field:   "",
		}
	}

	if err := v.validator.Struct(obj); err != nil {
		return v.formatValidationErrors(err)
	}

	return nil
}

// ToMap converts a struct to map[string]any.
// This is useful for converting model structs to maps for JSON responses
// or for further processing.
func (v *ModelValidator[T]) ToMap(obj *T) (map[string]any, error) {
	if obj == nil {
		return nil, fmt.Errorf("struct cannot be nil")
	}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	return result, nil
}

// ToJSON converts a struct to JSON bytes.
// This is useful for serializing structs to JSON format.
func (v *ModelValidator[T]) ToJSON(obj *T) ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("struct cannot be nil")
	}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	return jsonData, nil
}

// MapToStruct converts a map[string]any to a struct without validation.
// Use ValidateFromMap if you also need validation.
func (v *ModelValidator[T]) MapToStruct(data map[string]any) (*T, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	return v.JSONToStruct(jsonData)
}

// JSONToStruct converts JSON bytes to a struct without validation.
// Use ValidateFromJSON if you also need validation.
func (v *ModelValidator[T]) JSONToStruct(data []byte) (*T, error) {
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &obj, nil
}

// formatValidationErrors converts validator.ValidationErrors to our ValidationError format.
// It aggregates all field errors into a single error with detailed information.
func (v *ModelValidator[T]) formatValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationError{
			Message: err.Error(),
			Field:   "",
		}
	}

	// If there's only one error, return a simple ValidationError
	if len(validationErrors) == 1 {
		fieldErr := validationErrors[0]
		return &ValidationError{
			Message: v.getErrorMessage(fieldErr),
			Field:   fieldErr.Field(),
			Tag:     fieldErr.Tag(),
		}
	}

	// For multiple errors, create a ValidationErrors collection
	errors := make(ValidationErrors, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		errors = append(errors, &ValidationError{
			Message: v.getErrorMessage(fieldErr),
			Field:   fieldErr.Field(),
			Tag:     fieldErr.Tag(),
		})
	}

	return errors
}

// getErrorMessage generates a human-readable error message from a validator.FieldError.
func (v *ModelValidator[T]) getErrorMessage(fieldErr validator.FieldError) string {
	fieldName := fieldErr.Field()
	tag := fieldErr.Tag()
	param := fieldErr.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fieldName, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fieldName, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fieldName, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fieldName, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fieldName, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fieldName, param)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fieldName)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fieldName)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", fieldName)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", fieldName)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", fieldName)
	default:
		// For unknown tags, provide a generic message
		if param != "" {
			return fmt.Sprintf("%s failed validation for tag '%s' with parameter '%s'", fieldName, tag, param)
		}
		return fmt.Sprintf("%s failed validation for tag '%s'", fieldName, tag)
	}
}

// GetValidator returns the underlying validator instance for advanced usage.
// This allows custom validation rules to be registered if needed.
func (v *ModelValidator[T]) GetValidator() *validator.Validate {
	return v.validator
}

// ValidateWithCustomRules validates a struct and allows custom validation functions.
// This is useful when you need to add context-specific validation beyond struct tags.
func (v *ModelValidator[T]) ValidateWithCustomRules(obj *T, customRules func(*T) error) error {
	// First validate struct tags
	if err := v.ValidateStruct(obj); err != nil {
		return err
	}

	// Then apply custom rules if provided
	if customRules != nil {
		if err := customRules(obj); err != nil {
			return err
		}
	}

	return nil
}

// IsValid checks if a struct is valid without returning detailed errors.
// Returns true if valid, false otherwise.
func (v *ModelValidator[T]) IsValid(obj *T) bool {
	return v.ValidateStruct(obj) == nil
}

// ValidateAndConvert validates input and converts it to the target type.
// This is a convenience method that combines validation and conversion.
// Input can be []byte, map[string]any, or *T.
func (v *ModelValidator[T]) ValidateAndConvert(input any) (*T, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	switch val := input.(type) {
	case []byte:
		return v.ValidateFromJSON(val)
	case map[string]any:
		return v.ValidateFromMap(val)
	case *T:
		if err := v.ValidateStruct(val); err != nil {
			return nil, err
		}
		return val, nil
	case T:
		ptr := &val
		if err := v.ValidateStruct(ptr); err != nil {
			return nil, err
		}
		return ptr, nil
	default:
		// Try to convert via JSON marshaling/unmarshaling
		jsonData, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("unsupported input type: %s, failed to marshal: %w", reflect.TypeOf(input).String(), err)
		}
		return v.ValidateFromJSON(jsonData)
	}
}
