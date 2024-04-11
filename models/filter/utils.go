package filter

import "encoding/json"

func applyIfNotNil[T any](conditions map[string]interface{}, columnName string, value *T) {
	if value != nil {
		conditions[columnName] = *value
	}
}

func applyIfNotEmptySlice[T any](conditions map[string]interface{}, columnName string, value []T) {
	if len(value) > 0 {
		conditions[columnName] = value
	}
}

func applyConditionsIfNotNil(conditions map[string]interface{}, columnName string, nestedConditions map[string]interface{}) {
	if len(nestedConditions) > 0 {
		conditions[columnName] = nestedConditions
	}
}

func ptr[T any](value T) *T {
	return &value
}

func fromJSON[T any](raw string) T {
	var filter T
	_ = json.Unmarshal([]byte(raw), &filter)
	return filter
}
