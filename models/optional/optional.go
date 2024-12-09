package optional

// Param is a pointer to a value of type T.
type Param[T any] *T

// Of returns a pointer to a value of type T.
func Of[T any](v T) Param[T] {
	return &v
}
