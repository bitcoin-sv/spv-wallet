package chainstate

import (
	"context"
	"strings"
	"testing"
	"time"

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
		res := c.Broadcast(
			context.Background(), broadcastExample1TxID, broadcastExample1TxHex, RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.NotNil(t, res)
		require.Nil(t, res.Failure)

		miners := strings.Split(res.Provider, ",")
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
			WithBroadcastClient(bc),
		)

		// when
		res := c.Broadcast(
			context.Background(), broadcastExample1TxID, broadcastExample1TxHex, RawTx, defaultBroadcastTimeOut,
		)

		// then
		require.NotNil(t, res)
		require.Nil(t, res.Failure)

		assert.Equal(t, ProviderBroadcastClient, res.Provider)
	})

	t.Run("broadcast - success (multiple broadcast-client)", func(t *testing.T) {
		// given
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockFailure).
			WithMockArc(broadcast_client_mock.MockFailure).
			WithMockArc(broadcast_client_mock.MockSuccess).
			WithMockArc(broadcast_client_mock.MockFailure).
			Build()
		c := NewTestClient(
			context.Background(), t,
			WithBroadcastClient(bc),
		)

		// when
		res := c.Broadcast(
			context.Background(), broadcastExample1TxID, broadcastExample1TxHex, RawTx, 1*time.Second,
		)

		// then
		require.NotNil(t, res)
		require.Nil(t, res.Failure)

		assert.Equal(t, ProviderBroadcastClient, res.Provider)
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
