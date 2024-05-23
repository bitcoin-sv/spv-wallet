package datastore

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testClient will generate a test client
func testClient(ctx context.Context, t *testing.T, opts ...ClientOps) (ClientInterface, func()) {
	client, err := NewClient(ctx, opts...)
	require.NoError(t, err)
	require.NotNil(t, client)
	return client, func() {
		_ = client.Close(ctx)
	}
}

// TestClient_IsDebug will test the method IsDebug()
func TestClient_IsDebug(t *testing.T) {
	t.Run("toggle debug", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithDebugging())
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

// TestClient_Debug will test the method Debug()
func TestClient_Debug(t *testing.T) {
	t.Run("turn debug on", func(t *testing.T) {
		c, err := NewClient(context.Background())
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

// TestClient_DebugLog will test the method DebugLog()
func TestClient_DebugLog(t *testing.T) {
	t.Run("write debug log", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithDebugging())
		require.NotNil(t, c)
		require.NoError(t, err)

		c.DebugLog(context.Background(), "test message")
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

// TestClient_Engine will test the method Engine()
func TestClient_Engine(t *testing.T) {
	t.Run("[sqlite] - get engine", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithSQLite(&SQLiteConfig{
			DatabasePath: "",
			Shared:       true,
		}))
		assert.NotNil(t, c)
		require.NoError(t, err)
		assert.Equal(t, SQLite, c.Engine())
	})

	t.Run("[mongo] - failed to load", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithMongo(&MongoDBConfig{
			DatabaseName: "test",
			Transactions: false,
			URI:          "",
		}))
		assert.Nil(t, c)
		require.Error(t, err)
	})

	// todo: Postgresql
}

// TestClient_GetTableName will test the method GetTableName()
func TestClient_GetTableName(t *testing.T) {
	t.Run("table prefix", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithDebugging(), WithSQLite(&SQLiteConfig{
			CommonConfig: CommonConfig{
				TablePrefix: testTablePrefix,
			},
			DatabasePath: "",
			Shared:       true,
		}))
		require.NotNil(t, c)
		require.NoError(t, err)

		tableName := c.GetTableName(testModelName)
		assert.Equal(t, testTablePrefix+"_"+testModelName, tableName)
	})

	t.Run("no table prefix", func(t *testing.T) {
		c, err := NewClient(context.Background(), WithDebugging(), WithSQLite(&SQLiteConfig{
			CommonConfig: CommonConfig{
				TablePrefix: "",
			},
			DatabasePath: "",
			Shared:       true,
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
