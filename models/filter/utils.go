package filter

import "reflect"

func applyIfNotNil[T any](conditions map[string]interface{}, columnName string, value *T) {
	if value != nil {
		conditions[columnName] = &value
	}
}

func applyIfNotNilFunc[T any](conditions map[string]interface{}, columnName string, value *T, transformer func(*T) interface{}) {
	if value != nil {
		transformed := transformer(value)
		if !reflect.ValueOf(transformed).IsNil() {
			conditions[columnName] = transformed
		}
	}
}

func ptr[T any](value T) *T {
	return &value
}
