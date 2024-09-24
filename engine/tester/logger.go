package tester

import (
	"testing"

	"github.com/rs/zerolog"
)

// Logger returns a logger that can be used as a dependency in tests.
func Logger(t testing.TB) zerolog.Logger {
	return zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)))
}

// NopLogger returns a logger that does nothing.
// Deprecated: Use Logger instead.
func NopLogger() zerolog.Logger {
	return zerolog.Nop()
}
