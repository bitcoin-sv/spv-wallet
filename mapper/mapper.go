package mapper

// MapWithoutIndex allows to pass a mapper function that doesn't rquire index
// to "github.com/samber/lo" Map() function
func MapWithoutIndex[T, R any](mapper func(item T) R) func(T, int) R {
	return func(item T, _ int) R {
		return mapper(item)
	}
}

// MapWithoutIndexWithError allows to pass a mapper function that doesn't rquire index and also returns an error
// "github.com/samber/lo" Map() function
func MapWithoutIndexWithError[T, R any](mapper func(item T) (R, error)) func(T, int) (R, error) {
	return func(item T, _ int) (R, error) {
		return mapper(item)
	}
}
