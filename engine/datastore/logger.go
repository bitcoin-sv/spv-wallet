package datastore

import (
	zLogger "github.com/mrz1836/go-logger"
	gLogger "gorm.io/gorm/logger"
)

// DatabaseLogWrapper is a special wrapper for the GORM logger
type DatabaseLogWrapper struct {
	zLogger.GormLoggerInterface
}

// LogMode will set the log level/mode
func (d *DatabaseLogWrapper) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newLogger := *d
	if level == gLogger.Info {
		newLogger.SetMode(zLogger.Info)
	} else if level == gLogger.Warn {
		newLogger.SetMode(zLogger.Warn)
	} else if level == gLogger.Error {
		newLogger.SetMode(zLogger.Error)
	} else if level == gLogger.Silent {
		newLogger.SetMode(zLogger.Silent)
	}

	return &newLogger
}
