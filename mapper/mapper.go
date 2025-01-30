package mapper

func MapWithoutIndex[T, R any](mapper func(item T) R) func(T, int) R {
	return func(item T, _ int) R {
		return mapper(item)
	}
}

func MapWithoutIndexWithError[T, R any](mapper func(item T) (R, error)) func(T, int) (R, error) {
	return func(item T, _ int) (R, error) {
		return mapper(item)
	}
}
