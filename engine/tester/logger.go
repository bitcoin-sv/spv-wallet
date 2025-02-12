package tester

import (
	"testing"

	"github.com/rs/zerolog"
)

// Logger returns a logger that can be used as a dependency in tests.
func Logger(t testing.TB) zerolog.Logger {
	logger := zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)))
	return logger.Level(zerolog.TraceLevel)
}
