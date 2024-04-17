package chainstate

import (
	"context"
	"net/http"
	"testing"
	"time"

	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-minercraft/v2"
)

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("basic defaults", func(t *testing.T) {
		c, err := NewClient(
			context.Background(),
			WithMinercraft(&MinerCraftBase{}),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, false, c.IsDebug())
		assert.Equal(t, MainNet, c.Network())
		assert.Nil(t, c.HTTPClient())
		assert.NotNil(t, c.Minercraft())
	})

	t.Run("custom http client", func(t *testing.T) {
		customClient := &http.Client{}
		c, err := NewClient(
			context.Background(),
			WithHTTPClient(customClient),
			WithMinercraft(&MinerCraftBase{}),
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
			WithMinercraft(&MinerCraftBase{}),
			WithBroadcastClient(customClient),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.NotNil(t, c.BroadcastClient())
		assert.Equal(t, customClient, c.BroadcastClient())
	})

	t.Run("custom minercraft client", func(t *testing.T) {
		customClient, err := minercraft.NewClient(
			minercraft.DefaultClientOptions(), &http.Client{}, minercraft.MAPI, nil, nil,
		)
		require.NoError(t, err)
		require.NotNil(t, customClient)

		var c ClientInterface
		c, err = NewClient(
			context.Background(),
			WithMinercraft(customClient),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.NotNil(t, c.Minercraft())
		assert.Equal(t, customClient, c.Minercraft())
	})

	t.Run("custom query timeout", func(t *testing.T) {
		timeout := 55 * time.Second
		c, err := NewClient(
			context.Background(),
			WithQueryTimeout(timeout),
			WithMinercraft(&MinerCraftBase{}),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, timeout, c.QueryTimeout())
	})

	t.Run("custom network - test", func(t *testing.T) {
		c, err := NewClient(
			context.Background(),
			WithNetwork(TestNet),
			WithMinercraft(&MinerCraftBase{}),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, TestNet, c.Network())
	})

	t.Run("no provider when using minercraft with customNet", func(t *testing.T) {
		_, err := NewClient(
			context.Background(),
			WithNetwork(StressTestNet),
			WithMinercraft(&MinerCraftBase{}),
			WithFeeUnit(MockDefaultFee),
		)
		require.Error(t, err)
	})

	t.Run("unreachable miners", func(t *testing.T) {
		_, err := NewClient(
			context.Background(),
			WithMinercraft(&minerCraftUnreachable{}),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingBroadcastMiners)
	})
}
