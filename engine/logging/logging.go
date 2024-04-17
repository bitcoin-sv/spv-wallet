package logging

import (
	"os"

	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog"
)

// GetDefaultLogger generates and returns a default logger instance.
func GetDefaultLogger() *zerolog.Logger {
	logger := ecszerolog.New(os.Stdout, ecszerolog.Level(zerolog.InfoLevel)).
		With().
		Timestamp().
		Caller().
		Str("application", "spv-wallet-default").
		Logger()

	return &logger
}
