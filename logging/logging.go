package logging

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog"
)

const (
	consoleLogFormat = "console"
	jsonLogFormat    = "json"
)

// CreateLogger create and configure zerolog logger based on app config.
func CreateLogger(instanceName, format, level string, logOrigin bool) (*zerolog.Logger, error) {
	var writer io.Writer
	if format == consoleLogFormat {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05.000",
		}
	} else {
		writer = os.Stdout
	}

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		err = errors.New("failed to parse log level: " + err.Error())
		return nil, err
	}

	logLevel := ecszerolog.Level(parsedLevel)
	origin := ecszerolog.Origin()
	var logger zerolog.Logger

	if logOrigin {
		logger = ecszerolog.New(writer, logLevel, origin).
			With().
			Str("application", instanceName).
			Logger()
	} else {
		logger = ecszerolog.New(writer, logLevel).
			With().
			Str("application", instanceName).
			Logger()
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local) //nolint:gosmopolitan // We want local time inside logger.
	}

	return &logger, nil
}

// GetDefaultLogger create and configure default zerolog logger. It should be used before config is loaded
func GetDefaultLogger() *zerolog.Logger {
	logger := ecszerolog.New(os.Stdout, ecszerolog.Level(zerolog.DebugLevel)).
		With().
		Caller().
		Str("application", "spv-wallet-default").
		Logger()

	return &logger
}
