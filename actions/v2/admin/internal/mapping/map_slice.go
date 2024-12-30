package mapping

func mapSlice[Source any, Output any](itemParser func(source *Source) Output, source []*Source) []Output {
	result := make([]Output, 0, len(source))
	for _, item := range source {
		result = append(result, itemParser(item))
	}
	return result
}
