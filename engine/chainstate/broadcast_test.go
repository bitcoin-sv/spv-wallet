package chainstate

import (
	"context"
	"strings"
	"testing"

	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-minercraft/v2"
)

func Test_doesErrorContain(t *testing.T) {
	t.Run("valid contains", func(t *testing.T) {
		success := doesErrorContain("this is the test message", []string{"another", "test message"})
		assert.Equal(t, true, success)
	})

	t.Run("valid contains - equal case", func(t *testing.T) {
		success := doesErrorContain("this is the TEST message", []string{"another", "test message"})
		assert.Equal(t, true, success)
	})

	t.Run("does not contain", func(t *testing.T) {
		success := doesErrorContain("this is the test message", []string{"another", "nope"})
		assert.Equal(t, false, success)
	})
}

// TestClient_Broadcast will test the method Broadcast()
func TestClient_Broadcast(t *testing.T) {
	t.Parallel()

	t.Run("error - missing tx id", func(t *testing.T) {
		// given
		c := NewTestClient(context.Background(), t,
			WithMinercraft(&MinerCraftBase{}))

		// when
		provider, err := c.Broadcast(
			context.Background(), "", onChainExample1TxHex, RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTransactionID)
		assert.Empty(t, provider)
	})

	t.Run("error - missing tx hex", func(t *testing.T) {
		// given
		c := NewTestClient(context.Background(), t,
			WithMinercraft(&MinerCraftBase{}))

		// when
		provider, err := c.Broadcast(
			context.Background(), onChainExample1TxID, "", RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTransactionHex)
		assert.Empty(t, provider)
	})
}

// TestClient_Broadcast_MAPI will test the method Broadcast() with MAPI
func TestClient_Broadcast_MAPI(t *testing.T) {
	t.Parallel()

	t.Run("broadcast - success (mAPI)", func(t *testing.T) {
		// given
		c := NewTestClient(
			context.Background(), t,
			WithMinercraft(&minerCraftBroadcastSuccess{}),
		)

		// when
		providers, err := c.Broadcast(
			context.Background(), broadcastExample1TxID, broadcastExample1TxHex, RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.NoError(t, err)
		miners := strings.Split(providers, ",")
		assert.GreaterOrEqual(t, len(miners), 1)
		assert.True(t, containsAtLeastOneElement(
			miners,
			minercraft.MinerTaal,
			minercraft.MinerMempool,
			minercraft.MinerGorillaPool,
			minercraft.MinerMatterpool,
		))
	})
}

// TestClient_Broadcast_BroadcastClient will test the method Broadcast() with BroadcastClient
func TestClient_Broadcast_BroadcastClient(t *testing.T) {
	t.Parallel()

	t.Run("broadcast - success (broadcast-client)", func(t *testing.T) {
		// given
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockSuccess).
			Build()
		c := NewTestClient(
			context.Background(), t,
			WithMinercraft(&MinerCraftBase{}),
			WithBroadcastClient(bc),
		)

		// when
		providers, err := c.Broadcast(
			context.Background(), broadcastExample1TxID, broadcastExample1TxHex, RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.NoError(t, err)
		miners := strings.Split(providers, ",")
		assert.GreaterOrEqual(t, len(miners), 1)
		assert.True(t, containsAtLeastOneElement(
			miners,
			ProviderBroadcastClient,
		))
	})
}

func containsAtLeastOneElement(coll1 []string, coll2 ...string) bool {
	m := make(map[string]bool)

	for _, element := range coll1 {
		m[element] = true
	}

	// Check if any element from bool  is present in the set
	for _, element := range coll2 {
		if m[element] {
			return true
		}
	}

	return false
}
