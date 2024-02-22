package tester

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestLoadMockRedis will test the method LoadMockRedis()
func TestLoadMockRedis(t *testing.T) {
	t.Run("valid mock redis", func(t *testing.T) {
		idleTimeout := 10 * time.Second
		maxConnTime := idleTimeout
		maxActive := 10
		maxIdle := maxActive
		client, conn := LoadMockRedis(idleTimeout, maxConnTime, maxActive, maxIdle)
		require.NotNil(t, client)
		require.NotNil(t, conn)
		require.NotNil(t, client.Pool.Get())
	})
}

// TestLoadRealRedis will test the method LoadRealRedis()
func TestLoadRealRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test: redis is required")
	}

	t.Run("valid real redis", func(t *testing.T) {
		idleTimeout := 10 * time.Second
		maxConnTime := idleTimeout
		maxActive := 10
		maxIdle := maxActive
		client, conn, err := LoadRealRedis(
			testRedisConnection, idleTimeout, maxConnTime,
			maxActive, maxIdle, false, false,
		)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.NotNil(t, conn)
		require.NotNil(t, client.Pool.Get())
	})
}
