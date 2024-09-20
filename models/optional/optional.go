package optional

type Param[T any] *T

func Of[T any](v T) Param[T] {
	return &v
}
