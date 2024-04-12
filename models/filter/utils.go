package filter

func applyIfNotNil[T any](conditions map[string]interface{}, columnName string, value *T) {
	if value != nil {
		conditions[columnName] = &value
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
