package datastore

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_IsDebug(t *testing.T) {
	t.Run("toggle debug", func(t *testing.T) {
		c, err := NewClient(WithDebugging())
		require.NotNil(t, c)
		require.NoError(t, err)

		assert.True(t, c.IsDebug())

		c.Debug(false)

		assert.False(t, c.IsDebug())
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

func TestClient_Debug(t *testing.T) {
	t.Run("turn debug on", func(t *testing.T) {
		c, err := NewClient()
		require.NotNil(t, c)
		require.NoError(t, err)

		assert.False(t, c.IsDebug())

		c.Debug(true)

		assert.True(t, c.IsDebug())
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

func TestClient_DebugLog(t *testing.T) {
	t.Run("write debug log", func(t *testing.T) {
		c, err := NewClient(WithDebugging())
		require.NotNil(t, c)
		require.NoError(t, err)

		c.DebugLog(context.Background(), "test message")
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

func TestClient_Engine(t *testing.T) {
	t.Run("[sqlite] - get engine", func(t *testing.T) {
		c, err := NewClient(WithSQLite(&SQLiteConfig{
			DatabasePath: "",
			Shared:       false,
		}))
		assert.NotNil(t, c)
		require.NoError(t, err)
		assert.Equal(t, SQLite, c.Engine())
	})

	// todo: Postgresql
}

func TestClient_GetTableName(t *testing.T) {
	t.Run("table prefix", func(t *testing.T) {
		c, err := NewClient(WithDebugging(), WithSQLite(&SQLiteConfig{
			CommonConfig: CommonConfig{
				TablePrefix: testTablePrefix,
			},
			DatabasePath: "",
			Shared:       false,
		}))
		require.NotNil(t, c)
		require.NoError(t, err)

		tableName := c.GetTableName(testModelName)
		assert.Equal(t, testTablePrefix+"_"+testModelName, tableName)
	})

	t.Run("no table prefix", func(t *testing.T) {
		c, err := NewClient(WithDebugging(), WithSQLite(&SQLiteConfig{
			CommonConfig: CommonConfig{
				TablePrefix: "",
			},
			DatabasePath: "",
			Shared:       false,
		}))
		require.NotNil(t, c)
		require.NoError(t, err)

		tableName := c.GetTableName(testModelName)
		assert.Equal(t, testModelName, tableName)
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}
