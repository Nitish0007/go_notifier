package helpers

import (
	"fmt"
	"reflect"
)

func TransformTo[T any](source any, target *T) error {
	sourceValue := reflect.ValueOf(source)
	targetValue := reflect.ValueOf(target)

	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}

	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target is not a pointer")
	}

	if targetValue.IsNil() {
		return fmt.Errorf("target is nil")
	}

	targetValue = targetValue.Elem()

	return nil
}