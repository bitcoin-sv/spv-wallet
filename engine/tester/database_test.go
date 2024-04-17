package tester

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSQLiteTestConfig will test the method SQLiteTestConfig()
func TestSQLiteTestConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid config", func(t *testing.T) {
		config := SQLiteTestConfig(true, true)
		require.NotNil(t, config)

		assert.Equal(t, true, config.Debug)
		assert.Equal(t, true, config.Shared)
		assert.Equal(t, 1, config.MaxIdleConnections)
		assert.Equal(t, 1, config.MaxOpenConnections)
		assert.NotEmpty(t, config.TablePrefix)
		assert.Empty(t, config.DatabasePath)
	})

	t.Run("no debug or sharing", func(t *testing.T) {
		config := SQLiteTestConfig(false, false)
		require.NotNil(t, config)

		assert.Equal(t, false, config.Debug)
		assert.Equal(t, false, config.Shared)
	})
}
