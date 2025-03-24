package taskmanager

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestWithTaskQ(t *testing.T) {
	t.Run("check type", func(t *testing.T) {
		opt := WithTaskqConfig(nil)
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("test applying nil config", func(t *testing.T) {
		options := &options{
			taskq: &taskqOptions{
				config: nil,
				queue:  nil,
				tasks:  nil,
			},
		}
		opt := WithTaskqConfig(nil)
		opt(options)
		assert.Nil(t, options.taskq.config)
	})

	t.Run("test applying valid config", func(t *testing.T) {
		options := &options{
			taskq: &taskqOptions{},
		}
		opt := WithTaskqConfig(DefaultTaskQConfig("test-queue"))
		opt(options)
		assert.NotNil(t, options.taskq.config)
	})
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithLogger(nil)
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		options := &options{}
		opt := WithLogger(nil)
		opt(options)
		assert.Nil(t, options.logger)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &options{}
		customLogger := zerolog.Nop()
		opt := WithLogger(&customLogger)
		opt(options)
		assert.Equal(t, &customLogger, options.logger)
	})
}
