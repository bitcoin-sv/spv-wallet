package filter

import "encoding/json"

func fromJSON[T any](raw string) T {
	var filter T
	err := json.Unmarshal([]byte(raw), &filter)
	if err != nil {
		panic("Cannot unmarshall " + err.Error())
	}
	return filter
}

func ptr[T any](value T) *T {
	return &value
}
