package taskmanager

import (
	"github.com/rs/zerolog"
	"github.com/vmihailenco/taskq/v3"
)

// Options allow functional options to be supplied
type Options func(c *options)

// WithTaskqConfig will set the taskq custom config
func WithTaskqConfig(config *taskq.QueueOptions) Options {
	return func(c *options) {
		if config != nil {
			c.taskq.config = config
		}
	}
}

// WithLogger will set the custom logger interface
func WithLogger(customLogger *zerolog.Logger) Options {
	return func(c *options) {
		if customLogger != nil {
			c.logger = customLogger
		}
	}
}
