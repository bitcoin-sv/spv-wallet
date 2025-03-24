package tester

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
