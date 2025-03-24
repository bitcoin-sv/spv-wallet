package datastore

import (
	"os"
	"testing"

	zLogger "github.com/mrz1836/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientOptions(t *testing.T) {
	t.Run("ensure default values", func(t *testing.T) {
		defaults := defaultClientOptions()
		require.NotNil(t, defaults)
		assert.Equal(t, Empty, defaults.engine)
		assert.NotNil(t, defaults.sqLite)
	})
}

func TestWithDebugging(t *testing.T) {
	t.Run("get opts", func(t *testing.T) {
		opt := WithDebugging()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opts := []ClientOps{WithDebugging()}
		c, err := NewClient(opts...)
		require.NotNil(t, c)
		require.NoError(t, err)

		assert.True(t, c.IsDebug())
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

func TestWithSQLite(t *testing.T) {
	t.Run("check type", func(t *testing.T) {
		opt := WithSQLite(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		options := &clientOptions{}
		opt := WithSQLite(nil)
		opt(options)
		assert.Equal(t, Engine(""), options.engine)
		assert.Nil(t, options.sqLite)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{}
		config := &SQLiteConfig{
			CommonConfig: CommonConfig{
				Debug:              true,
				MaxIdleConnections: 1,
				MaxOpenConnections: 1,
				TablePrefix:        "test",
			},
			DatabasePath:       "",
			ExistingConnection: nil,
			Shared:             false,
		}
		opt := WithSQLite(config)
		opt(options)
		assert.Equal(t, config, options.sqLite)
		assert.Equal(t, maxIdleConnectionsSQLite, options.sqLite.MaxIdleConnections)
		assert.Equal(t, SQLite, options.engine)
		assert.Equal(t, config.TablePrefix, options.tablePrefix)
		assert.True(t, options.debug)
	})
}

func TestWithSQL(t *testing.T) {
	t.Run("check type", func(t *testing.T) {
		opt := WithSQL("", nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty engine", func(t *testing.T) {
		options := &clientOptions{}
		opt := WithSQL(Empty, nil)
		opt(options)
		assert.Equal(t, Engine(""), options.engine)
		assert.Nil(t, options.sqlConfigs)
	})
}

func TestWithSQLConnection(t *testing.T) {
	t.Run("check type", func(t *testing.T) {
		opt := WithSQLConnection("", nil, testTablePrefix)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty engine", func(t *testing.T) {
		options := &clientOptions{}
		opt := WithSQLConnection(Empty, nil, testTablePrefix)
		opt(options)
		assert.Equal(t, Engine(""), options.engine)
		assert.Nil(t, options.sqlConfigs)
	})
}

func TestWithLogger(t *testing.T) {
	t.Run("check type", func(t *testing.T) {
		opt := WithLogger(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		options := &clientOptions{}
		opt := WithLogger(nil)
		opt(options)
		assert.Nil(t, options.logger)
	})

	t.Run("test applying valid logger", func(t *testing.T) {
		options := &clientOptions{}
		l := zLogger.NewGormLogger(true, 4)
		opt := WithLogger(l)
		opt(options)
		assert.Equal(t, l, options.logger)
	})
}
