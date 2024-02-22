package engine

import (
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
