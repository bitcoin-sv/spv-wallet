package mapper

func MapWithoutIndex[Source, Output any](mapper func(item Source) Output) func(Source, int) Output {
	return func(item Source, _ int) Output {
		return mapper(item)
	}
}
