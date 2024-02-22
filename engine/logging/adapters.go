package logging

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mrz1836/go-logger"
	"github.com/rs/zerolog"
)

// GormLoggerAdapter is an adapter for integrating the GORM library with the zerolog logger.
type GormLoggerAdapter struct {
	Logger     *zerolog.Logger
	logLevel   logger.GormLogLevel
	stackLevel int
}

// Error logs an error message.
func (a *GormLoggerAdapter) Error(_ context.Context, msg string, parameters ...any) {
	event := prepareBasicEvent(a.Logger.Error(), parameters...)
	event.Msgf(msg)
}

// Info logs an informative message.
func (a *GormLoggerAdapter) Info(_ context.Context, msg string, parameters ...any) {
	event := prepareBasicEvent(a.Logger.Info(), parameters...)
	event.Msgf(msg)
}

// Warn logs a warning message.
func (a *GormLoggerAdapter) Warn(_ context.Context, msg string, parameters ...any) {
	event := prepareBasicEvent(a.Logger.Warn(), parameters...)
	event.Msgf(msg)
}

// Trace logs the execution time and details of a database query.
// It checks the log level and the error status to determine the appropriate logging action.
// If the log level is set to logger.Error and there is an error (excluding "record not found"), it logs a warning message with the error details and the executed query.
// If the elapsed time exceeds the slow query threshold and the log level is set to logger.Warn, it logs a warning message with the slow query details.
// If the log level is set to logger.Info, it logs an executing query message with the query details.
func (a *GormLoggerAdapter) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if a.logLevel <= 1 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && a.logLevel >= logger.Error && (!strings.Contains(err.Error(), "record not found")):
		sql, rows := fc()
		event := prepareTraceEvent(a.Logger.Error(), elapsed, rows, sql)
		event.Str("error", err.Error()).
			Msg("warning executing query")
	case elapsed > logger.SlowQueryThreshold && a.logLevel >= logger.Warn:
		sql, rows := fc()
		event := prepareTraceEvent(a.Logger.Warn(), elapsed, rows, sql)
		event.Str("slow_log", fmt.Sprintf("SLOW SQL >= %v", logger.SlowQueryThreshold)).
			Msg("warning executing query")
	case a.logLevel == logger.Info:
		sql, rows := fc()
		event := prepareTraceEvent(a.Logger.Info(), elapsed, rows, sql)
		event.Msg("executing query")
	}
}

// GetStackLevel returns the stack level of the GormLoggerAdapter.
func (a *GormLoggerAdapter) GetStackLevel() int {
	return a.stackLevel
}

// SetStackLevel sets the stack level to the specified value.
func (a *GormLoggerAdapter) SetStackLevel(level int) {
	a.stackLevel = level
}

// GetMode returns the current log level of the GormLoggerAdapter.
func (a *GormLoggerAdapter) GetMode() logger.GormLogLevel {
	return a.logLevel
}

// SetMode sets the log level of the GormLoggerAdapter to the specified value.
// If the level is Silent, set the logLevel to 1.
// If the level is Error, set the logLevel to 2 and set the Logger level to zerolog.ErrorLevel.
// If the level is Warn, set the logLevel to 3 and set the Logger level to zerolog.WarnLevel.
// If the level is Info, set the logLevel to 4 and set the Logger level to zerolog.InfoLevel.
// There is no Debug mode inn GormLoggerAdapter.
func (a *GormLoggerAdapter) SetMode(level logger.GormLogLevel) logger.GormLoggerInterface {
	if level == logger.Silent {
		a.logLevel = 1
	} else if level == logger.Error {
		a.Logger.Level(zerolog.ErrorLevel)
		a.logLevel = 2
	} else if level == logger.Warn {
		a.Logger.Level(zerolog.WarnLevel)
		a.logLevel = 3
	} else if level == logger.Info {
		a.Logger.Level(zerolog.InfoLevel)
		a.logLevel = 4
	}

	return a
}

// CreateGormLoggerAdapter creates a new instance of the GormLoggerAdapter struct that wraps a zerolog.Logger and provides logging functionality for GORM.
// It takes a pointer to a zerolog.Logger and a serviceName string as parameters.
// It determines the log level based on the global log level of the zerolog.Logger and assigns the corresponding logger.GormLogLevel to the newly created GormLoggerAdapter.
// It creates a subservice logger by setting the "subservice" field in the zerolog.Logger context with the serviceName parameter and returns the GormLoggerAdapter.
func CreateGormLoggerAdapter(zLog *zerolog.Logger, serviceName string) *GormLoggerAdapter {
	level := zLog.GetLevel()
	var l logger.GormLogLevel

	if level == zerolog.ErrorLevel {
		l = logger.Error
	} else if level == zerolog.WarnLevel {
		l = logger.Warn
	} else if level == zerolog.InfoLevel {
		l = logger.Info
	} else if level == zerolog.DebugLevel {
		l = logger.Info
	} else {
		l = logger.Silent
	}

	subserviceLogger := zLog.With().Str("subservice", serviceName).Logger()

	return &GormLoggerAdapter{
		Logger:   &subserviceLogger,
		logLevel: l,
	}
}

// fileWithLineNum returns the file path and line number of the caller function.
func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasSuffix(file, "_test.go") && // Skip test files
			!strings.Contains(file, "gorm.go") && // This is our local "gorm.go" file
			!strings.Contains(file, "callbacks.go") && // This file is a helper for GORM
			!strings.Contains(file, "finisher_api.go")) { // This file is a helper for GORM
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

// prepareParameters is a function that takes a variadic number of parameters and returns a slice of logger.KeyValue.
// It iterates through the parameters and creates a logger.Parameter for each parameter, using the index as the key and the parameter value as the value.
// The created logger.Parameter is then appended to the keyValues slice.
func prepareParameters(params ...any) []logger.KeyValue {
	var keyValues []logger.KeyValue
	if len(params) > 0 {
		for index, val := range params {
			parameter := logger.Parameter{K: fmt.Sprintf("param_%d", index), V: val}
			keyValues = append(keyValues, &parameter)
		}
	}
	return keyValues
}

// prepareTraceEvent prepares a trace event.
// The function adds the following fields to the event:
// - "file": The file name and line number where the event occurred, obtained from the fileWithLineNum function.
// - "duration": The duration of the event in milliseconds, formatted as a string with three decimal places.
// - "rows": The number of rows affected by the event.
// - "sql": The SQL query associated with the event.
func prepareTraceEvent(event *zerolog.Event, elapsed time.Duration, rows int64, sql string) *zerolog.Event {
	return event.
		Str("file", fileWithLineNum()).
		Str("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)).
		Int64("rows", rows).
		Str("sql", sql)
}

// prepareBasicEvent prepares a basic event by adding parameters to the provided zerolog.Event.
// The parameters are converted to logger.KeyValue pairs using the prepareParameters function,
// and each key-value pair is added to the event using the event.Str(key, value) method.
func prepareBasicEvent(event *zerolog.Event, parameters ...any) *zerolog.Event {
	params := prepareParameters(parameters...)
	for _, keyValue := range params {
		event.Str(keyValue.Key(), fmt.Sprint(keyValue.Value()))
	}
	return event
}
