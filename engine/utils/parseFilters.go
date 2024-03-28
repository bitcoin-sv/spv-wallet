package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// FilterMapByStructFields returns a new map containing only the fields that exist in the provided struct.
func FilterMapByStructFields(fieldMap map[string]interface{}, filtersStruct interface{}) (map[string]interface{}, error) {
	// Check if filtersStruct is a pointer to a struct
	filtersValue := reflect.ValueOf(filtersStruct)
	if filtersValue.Kind() != reflect.Ptr || filtersValue.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("filtersStruct must be a pointer to a struct")
	}

	// Get the type of the struct
	structType := filtersValue.Elem().Type()

	// Create a new map to store filtered fields
	filteredMap := make(map[string]interface{})

	// Iterate over the keys of the field map
	for fieldName, fieldValue := range fieldMap {
		// Check if the field exists in the struct
		found := false
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			tag := field.Tag.Get("json")
			if strings.Contains(tag, fieldName) && fieldValue != nil {
				filteredMap[fieldName] = fieldValue
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("field '%s' does not exist in the struct", fieldName)
		}
	}
	return filteredMap, nil
}
