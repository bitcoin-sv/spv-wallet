package logging

import (
	"io"
	"os"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog"
)

const (
	consoleLogFormat = "console"
	jsonLogFormat    = "json"
)

// GetDefaultLogger create and configure default zerolog logger. It should be used before config is loaded
func GetDefaultLogger() zerolog.Logger {
	logger, err := createLogger("spv-wallet-default", jsonLogFormat, "debug", true)
	if err != nil {
		panic(err)
	}
	return logger
}

func CreateLoggerWithConfig(config *config.AppConfig) (zerolog.Logger, error) {
	loggingConfig := config.Logging
	return createLogger(loggingConfig.InstanceName, loggingConfig.Format, loggingConfig.Level, loggingConfig.LogOrigin)
}

func createLogger(instanceName, format, level string, logOrigin bool) (zerolog.Logger, error) {
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
		err = spverrors.Wrapf(err, "failed to parse log level")
		return zerolog.Nop(), err
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

	return logger, nil
}
