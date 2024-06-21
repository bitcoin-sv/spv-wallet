package chainstate

import (
	"context"
	"net/http"
	"testing"
	"time"

	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockSuccess).
		Build()

	t.Run("basic defaults", func(t *testing.T) {
		c, err := NewClient(
			context.Background(),
			WithBroadcastClient(bc),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, false, c.IsDebug())
		assert.Equal(t, MainNet, c.Network())
		assert.Nil(t, c.HTTPClient())
	})

	t.Run("custom http client", func(t *testing.T) {
		customClient := &http.Client{}
		c, err := NewClient(
			context.Background(),
			WithHTTPClient(customClient),
			WithBroadcastClient(bc),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.NotNil(t, c.HTTPClient())
		assert.Equal(t, customClient, c.HTTPClient())
	})

	t.Run("custom broadcast client", func(t *testing.T) {
		arcConfig := broadcast_client.ArcClientConfig{
			Token:  "",
			APIUrl: "https://arc.gorillapool.io",
		}
		logger := zerolog.Nop()
		customClient := broadcast_client.Builder().WithArc(arcConfig, &logger).Build()
		require.NotNil(t, customClient)
		c, err := NewClient(
			context.Background(),
			WithBroadcastClient(customClient),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.NotNil(t, c.BroadcastClient())
		assert.Equal(t, customClient, c.BroadcastClient())
	})

	t.Run("custom query timeout", func(t *testing.T) {
		timeout := 55 * time.Second
		c, err := NewClient(
			context.Background(),
			WithQueryTimeout(timeout),
			WithBroadcastClient(bc),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, timeout, c.QueryTimeout())
	})

	t.Run("custom network - test", func(t *testing.T) {
		c, err := NewClient(
			context.Background(),
			WithNetwork(TestNet),
			WithBroadcastClient(bc),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, TestNet, c.Network())
	})
}
