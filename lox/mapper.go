package lox

// MappingFn allows to pass a mapper function that doesn't rquire index
// to "github.com/samber/lo" Map() function
func MappingFn[T, R any](mapper func(item T) R) func(T, int) R {
	return func(item T, _ int) R {
		return mapper(item)
	}
}

// MappingFnWithError allows to pass a mapper function that doesn't rquire index and also returns an error
// "github.com/samber/lo" Map() function
func MappingFnWithError[T, R any](mapper func(item T) (R, error)) func(T, int) (R, error) {
	return func(item T, _ int) (R, error) {
		return mapper(item)
	}
}

// MapAndCollect takes ErrorCollector and a function mapper that can return an error.
// If the mapper function returns an error it is joined to the errors field of ErrorCollector
func MapAndCollect[T, R any](catcher *ErrorCollector, iteratee func(item T, index int) (R, error)) func(item T, index int) R {
	return func(item T, index int) R {
		res, err := iteratee(item, index)
		if err != nil {
			catcher.Collect(err)
		}
		return res
	}
}
