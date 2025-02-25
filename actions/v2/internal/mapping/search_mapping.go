package mapping

import (
	"reflect"
	"strings"
)

// MapToDBConditions extracts non-nil fields from any struct, excluding pagination fields
func MapToDBConditions[T any](params T) map[string]interface{} {
	conditions := map[string]interface{}{}

	excludedFields := map[string]bool{
		"page":   true,
		"size":   true,
		"sort":   true,
		"sortBy": true,
	}

	val := reflect.ValueOf(params)
	typ := reflect.TypeOf(params)

	if typ.Kind() != reflect.Struct {
		return conditions
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		jsonTag := fieldType.Tag.Get("json")
		dbField := strings.Split(jsonTag, ",")[0]

		if dbField == "" {
			dbField = fieldType.Name
		}

		if _, found := excludedFields[dbField]; found {
			continue
		}

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			conditions[dbField] = field.Elem().Interface()
		}
	}

	return conditions
}

// GetPointerValue is a generic helper function to dereference pointers safely
func GetPointerValue[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}
