package utils

// MapSlice maps a slice of items to a slice of another type of items
func MapSlice[Source any, Output any](source []Source, itemParser func(source Source) Output) []Output {
	result := make([]Output, 0, len(source))
	for _, item := range source {
		result = append(result, itemParser(item))
	}
	return result
}
