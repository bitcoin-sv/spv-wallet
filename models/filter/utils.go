package filter

import (
	"errors"
	"reflect"
	"strings"
)

func applyIfNotNil[T any](conditions map[string]interface{}, columnName string, value *T) {
	if value != nil {
		conditions[columnName] = *value
	}
}

func applyConditionsIfNotNil(conditions map[string]interface{}, columnName string, nestedConditions map[string]interface{}) {
	if len(nestedConditions) > 0 {
		conditions[columnName] = nestedConditions
	}
}

// strOption checks (case-insensitive) if value is in options, if it is, it returns a pointer to the value, otherwise it returns an error
func checkStrOption(value string, options ...string) (string, error) {
	for _, opt := range options {
		if strings.EqualFold(value, opt) {
			return opt, nil
		}
	}
	return "", errors.New("invalid filter option")
}

func checkAndApplyStrOption(conditions map[string]interface{}, columnName string, value *string, options ...string) error {
	if value == nil {
		return nil
	}
	opt, err := checkStrOption(*value, options...)
	if err != nil {
		return err
	}
	conditions[columnName] = opt
	return nil
}

// getEnumValues gets the tag "enums" of a field by fieldName of a provided struct
func getEnumValues[T any](fieldName string) []string {
	t := reflect.TypeOf(*new(T))
	field, found := t.FieldByName(fieldName)
	if !found {
		return nil
	}
	enums := field.Tag.Get("enums")
	if enums == "" {
		return nil
	}
	options := strings.Split(enums, ",")
	for i := 0; i < len(options); i++ {
		options[i] = strings.TrimSpace(options[i])
	}
	return options
}
