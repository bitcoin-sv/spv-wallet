package chainstate

import (
	"context"
	"net/http"
	"testing"
	"time"

	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithNewRelic will test the method WithNewRelic()
func TestWithNewRelic(t *testing.T) {
	t.Run("get opts", func(t *testing.T) {
		opt := WithNewRelic()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockSuccess).
			Build()
		opts := []ClientOps{
			WithNewRelic(),
			WithBroadcastClient(bc),
		}
		logger := zerolog.Nop()
		c, err := NewClient(context.Background(), &logger, opts...)
		require.NotNil(t, c)
		require.NoError(t, err)

		assert.Equal(t, true, c.IsNewRelicEnabled())
	})
}

// TestWithHTTPClient will test the method WithHTTPClient()
func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithHTTPClient(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithHTTPClient(nil)
		opt(options)
		assert.Nil(t, options.config.httpClient)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		customClient := &http.Client{}
		opt := WithHTTPClient(customClient)
		opt(options)
		assert.Equal(t, customClient, options.config.httpClient)
	})
}

// TestWithBroadcastClient will test the method WithBroadcastClient()
func TestWithBroadcastClient(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithBroadcastClient(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithBroadcastClient(nil)
		opt(options)
		assert.Nil(t, options.config.broadcastClient)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		customClient := broadcast_client_mock.Builder().WithMockArc(broadcast_client_mock.MockSuccess).Build()
		opt := WithBroadcastClient(customClient)
		opt(options)
		assert.Equal(t, customClient, options.config.broadcastClient)
	})
}

// TestWithQueryTimeout will test the method WithQueryTimeout()
func TestWithQueryTimeout(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithQueryTimeout(0)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty value", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithQueryTimeout(0)
		opt(options)
		assert.Equal(t, time.Duration(0), options.config.queryTimeout)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithQueryTimeout(10 * time.Second)
		opt(options)
		assert.Equal(t, 10*time.Second, options.config.queryTimeout)
	})
}

// TestWithNetwork will test the method WithNetwork()
func TestWithNetwork(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithNetwork("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty string", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithNetwork("")
		opt(options)
		assert.Equal(t, Network(""), options.config.network)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithNetwork(TestNet)
		opt(options)
		assert.Equal(t, TestNet, options.config.network)
	})
}

// TestWithUserAgent will test the method WithUserAgent()
func TestWithUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithUserAgent("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty string", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithUserAgent("")
		opt(options)
		assert.Equal(t, "", options.userAgent)
	})

	t.Run("test applying option", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithUserAgent("test agent")
		opt(options)
		assert.Equal(t, "test agent", options.userAgent)
	})
}

// TestWithExcludedProviders will test the method WithExcludedProviders()
func TestWithExcludedProviders(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithExcludedProviders(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying empty string", func(t *testing.T) {
		options := &clientOptions{
			config: &syncConfig{},
		}
		opt := WithExcludedProviders([]string{""})
		opt(options)
		assert.Equal(t, []string{""}, options.config.excludedProviders)
	})
}
