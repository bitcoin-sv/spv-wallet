package lox

// Iteratee represents a function that processes an item of type T at a specific index,
// returning a result of type R. It is useful for operations requiring both the item and its index.
type Iteratee[T, R any] func(item T, index int) R

// NoIndexIteratee represents a function that processes an item of type T
// without requiring its index, returning a result of type R.
// It is useful for operations where the index is irrelevant.
type NoIndexIteratee[T, R any] func(item T) R

// IterateeWithError represents a function that processes an item of type T at a specific index,
// returning a result of type R and an error if the operation fails.
// It is useful for operations that might encounter errors during processing.
type IterateeWithError[T, R any] func(item T, index int) (R, error)

// NoIndexIterateeWithError represents a function that processes an item of type T
// without requiring its index, returning a result of type R and an error if the operation fails.
// It is useful for operations where the index is irrelevant and error handling is required.
type NoIndexIterateeWithError[T, R any] func(item T) (R, error)

// MappingFn allows to pass a mapper function that doesn't require index
// to "github.com/samber/lo" Map() function
func MappingFn[T, R any](mapper NoIndexIteratee[T, R]) Iteratee[T, R] {
	return func(item T, _ int) R {
		return mapper(item)
	}
}

// MappingFnWithError allows to pass a mapper function that doesn't require index and also returns an error
// "github.com/samber/lo" Map() function
func MappingFnWithError[T, R any](mapper NoIndexIterateeWithError[T, R]) IterateeWithError[T, R] {
	return func(item T, _ int) (R, error) {
		return mapper(item)
	}
}

// MapAndCollect takes ErrorCollector and a function mapper that can return an error.
// If the mapper function returns an error it is joined to the errors field of ErrorCollector
func MapAndCollect[T, R any](catcher *ErrorCollector, iteratee NoIndexIterateeWithError[T, R]) Iteratee[T, R] {
	return func(item T, _ int) R {
		res, err := iteratee(item)
		if err != nil {
			catcher.Collect(err)
		}
		return res
	}
}

// MapEntriesOrError manipulates a map entries and transforms it to a map of another type.
// If the iteratee function returns an error, the function returns the error.
func MapEntriesOrError[K1 comparable, V1 any, K2 comparable, V2 any](in map[K1]V1, iteratee func(key K1, value V1) (K2, V2, error)) (map[K2]V2, error) {
	result := make(map[K2]V2, len(in))

	for k1 := range in {
		k2, v2, err := iteratee(k1, in[k1])
		if err != nil {
			return nil, err
		}
		result[k2] = v2
	}

	return result, nil
}
