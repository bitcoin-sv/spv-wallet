package lox

import "github.com/samber/lo"

type IfElse[T any] interface {
	Else(result T) T
}

func Unwrap[T any](value *T) IfElse[T] {
	return lo.IfF(value != nil, func() T {
		return *value
	})
}
