package filter

import (
	"errors"
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
func strOption(value *string, options ...string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	for _, opt := range options {
		if strings.EqualFold(*value, opt) {
			s := string(opt)
			return &s, nil
		}
	}
	return nil, errors.New("Invalid option: " + *value)
}
