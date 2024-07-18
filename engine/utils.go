package engine

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog"
)

func recoverAndLog(log *zerolog.Logger) {
	if err := recover(); err != nil {
		log.Error().Msgf(
			"panic: %v - stack trace: %v", err,
			strings.ReplaceAll(string(debug.Stack()), "\n", ""),
		)
	}
}

// finds the first element in a collection that satisfies a specified condition.
func find[E any](collection []E, predicate func(E) bool) *E {
	for _, v := range collection {
		if predicate(v) {
			return &v
		}
	}
	return nil
}

func contains[E any](collection []E, predicate func(E) bool) bool {
	return find(collection, predicate) != nil
}

// EnumStringMapper is a helper struct that maps strings to "enum" like constants.
type EnumStringMapper[T fmt.Stringer] struct {
	elements map[string]T
}

// NewEnumStringMapper creates a new EnumStringMapper.
func NewEnumStringMapper[T fmt.Stringer](elements ...T) EnumStringMapper[T] {
	m := make(map[string]T)
	for _, element := range elements {
		m[element.String()] = element
	}
	return EnumStringMapper[T]{
		elements: m,
	}
}

// Get returns the "enum" value for the given string.
func (m *EnumStringMapper[T]) Get(key string) (T, bool) {
	value, ok := m.elements[key]
	return value, ok
}
