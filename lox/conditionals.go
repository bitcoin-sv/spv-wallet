package lox

import "github.com/samber/lo"

// IfElse is an interface that satisfies github.com/samber/lo IfF
// function return value
type IfElse[T any] interface {
	Else(result T) T
}

// Unwrap is a helper function that takes a pointer and returns its
// unwrapped value if it is not nil
func Unwrap[T any](value *T) IfElse[T] {
	return lo.IfF(value != nil, func() T {
		return *value
	})
}
