package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog"
)

const (
	consoleLogFormat = "console"
	jsonLogFormat    = "json"
)

// GetDefaultLogger create and configure default zerolog logger. It should be used before config is loaded
func GetDefaultLogger() zerolog.Logger {
	return CreateLogger(writerFor(jsonLogFormat), "spv-wallet-default", zerolog.DebugLevel, true)
}

// CreateLoggerWithConfig creates a logger based on the given config
func CreateLoggerWithConfig(config *config.AppConfig) (zerolog.Logger, error) {
	loggingConfig := config.Logging
	parsedLevel, err := zerolog.ParseLevel(loggingConfig.Level)
	if err != nil {
		return zerolog.Nop(), spverrors.Wrapf(err, "failed to parse log level")
	}

	return CreateLogger(writerFor(loggingConfig.Format), loggingConfig.InstanceName, parsedLevel, loggingConfig.LogOrigin), nil
}

// CreateLogger creates a logger with the given writer, instance name and log level
func CreateLogger(writer io.Writer, instanceName string, level zerolog.Level, logOrigin bool) zerolog.Logger {
	options := []ecszerolog.Option{
		ecszerolog.Level(level),
	}

	if logOrigin {
		options = append(options, ecszerolog.Origin())
	}

	logger := ecszerolog.New(writer, options...).
		With().
		Str("application", instanceName).
		Logger()

	// NOTE: zerolog.New() overwrites global handlers, so we need to set them AFTER creating the logger
	setGlobalHandlers()

	return logger
}

func setGlobalHandlers() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local) //nolint:gosmopolitan // We want local time inside logger.
	}

	zerolog.ErrorMarshalFunc = func(err error) any {
		if errorx.Cast(err) != nil {
			return fmt.Sprintf("%v", err)
		}
		return spverrors.UnfoldError(err)
	}

	zerolog.ErrorStackMarshaler = func(err error) any {
		if errorx.Cast(err) != nil {
			fullMessage := fmt.Sprintf("%+v", err)
			const startingNewLine = "\n "
			const stackTraceMarker = startingNewLine + "at "
			stackTraceStart := strings.Index(fullMessage, stackTraceMarker)
			if stackTraceStart == -1 {
				return nil
			}
			return formatStackTrace(fullMessage[stackTraceStart+len(startingNewLine):])
		}
		return nil
	}
}

func writerFor(format string) io.Writer {
	switch format {
	case consoleLogFormat:
		return zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05.000",
		}
	default:
		return os.Stdout
	}
}

func formatStackTrace(stackMsg string) []string {
	stackMsg = strings.ReplaceAll(stackMsg, "\n\t", " ")
	return strings.Split(stackMsg, "\n ")
}
