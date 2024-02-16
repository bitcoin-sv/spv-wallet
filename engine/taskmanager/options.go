package taskmanager

import (
	"github.com/rs/zerolog"
	taskq "github.com/vmihailenco/taskq/v3"
)

// TaskManagerOptions allow functional options to be supplied
type TaskManagerOptions func(c *options)

// WithNewRelic will enable the NewRelic wrapper
func WithNewRelic() TaskManagerOptions {
	return func(c *options) {
		c.newRelicEnabled = true
	}
}

// WithTaskqConfig will set the taskq custom config
func WithTaskqConfig(config *taskq.QueueOptions) TaskManagerOptions {
	return func(c *options) {
		if config != nil {
			c.taskq.config = config
		}
	}
}

// WithLogger will set the custom logger interface
func WithLogger(customLogger *zerolog.Logger) TaskManagerOptions {
	return func(c *options) {
		if customLogger != nil {
			c.logger = customLogger
		}
	}
}
