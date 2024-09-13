package tester

import "github.com/rs/zerolog"

// Logger returns a logger that can be used as a dependency in tests.
func Logger() zerolog.Logger {
	return zerolog.Nop()
}
