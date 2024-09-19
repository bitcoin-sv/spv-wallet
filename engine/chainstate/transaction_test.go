package chainstate

import (
	"context"
	"testing"

	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	broadcast_fixtures "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_Transaction will test the method QueryTransaction()
func TestClient_Transaction(t *testing.T) {
	t.Parallel()
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockSuccess).
		Build()

	t.Run("error - missing id", func(t *testing.T) {
		// given
		c := NewTestClient(context.Background(), t, WithBroadcastClient(bc))

		// when
		info, err := c.QueryTransaction(
			context.Background(), "", RequiredOnChain, defaultQueryTimeOut,
		)

		// then
		require.Error(t, err)
		require.Nil(t, info)
		assert.ErrorIs(t, err, spverrors.ErrInvalidTransactionID)
	})

	t.Run("error - missing requirements", func(t *testing.T) {
		// given
		c := NewTestClient(context.Background(), t, WithBroadcastClient(bc))

		// when
		info, err := c.QueryTransaction(
			context.Background(), onChainExample1TxID,
			"", defaultQueryTimeOut,
		)

		// then
		require.Error(t, err)
		require.Nil(t, info)
		assert.ErrorIs(t, err, spverrors.ErrInvalidRequirements)
	})
}

func TestClient_Transaction_BroadcastClient(t *testing.T) {
	t.Parallel()

	t.Run("query transaction success - broadcastClient", func(t *testing.T) {
		// given
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockSuccess).
			Build()
		c := NewTestClient(
			context.Background(), t,
			WithBroadcastClient(bc),
		)

		// when
		info, err := c.QueryTransaction(
			context.Background(), onChainExampleArcTxID,
			RequiredInMempool, defaultQueryTimeOut,
		)

		// then
		require.NoError(t, err)
		require.NotNil(t, info)
		assert.Equal(t, onChainExampleArcTxID, info.ID)
		assert.Equal(t, broadcast_fixtures.TxBlockHash, info.BlockHash)
		assert.Equal(t, broadcast_fixtures.TxBlockHeight, info.BlockHeight)
		assert.Equal(t, broadcast_fixtures.ProviderMain, info.Provider)
	})

	t.Run("valid - stress test network - broadcastClient", func(t *testing.T) {
		// given
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockSuccess).
			Build()
		c := NewTestClient(
			context.Background(), t,
			WithBroadcastClient(bc),
			WithNetwork(StressTestNet),
		)

		// when
		info, err := c.QueryTransaction(
			context.Background(), onChainExampleArcTxID,
			RequiredInMempool, defaultQueryTimeOut,
		)

		// then
		require.NoError(t, err)
		require.NotNil(t, info)
		assert.Equal(t, onChainExampleArcTxID, info.ID)
		assert.Equal(t, broadcast_fixtures.TxBlockHash, info.BlockHash)
		assert.Equal(t, broadcast_fixtures.TxBlockHeight, info.BlockHeight)
		assert.Equal(t, broadcast_fixtures.ProviderMain, info.Provider)
	})

	t.Run("valid - test network - broadcast", func(t *testing.T) {
		// given
		bc := broadcast_client_mock.Builder().
			WithMockArc(broadcast_client_mock.MockSuccess).
			Build()
		c := NewTestClient(
			context.Background(), t,
			WithBroadcastClient(bc),
			WithNetwork(TestNet),
		)

		// when
		info, err := c.QueryTransaction(
			context.Background(), onChainExampleArcTxID,
			RequiredInMempool, defaultQueryTimeOut,
		)

		// then
		require.NoError(t, err)
		require.NotNil(t, info)
		assert.Equal(t, onChainExampleArcTxID, info.ID)
		assert.Equal(t, broadcast_fixtures.TxBlockHash, info.BlockHash)
		assert.Equal(t, broadcast_fixtures.TxBlockHeight, info.BlockHeight)
		assert.Equal(t, broadcast_fixtures.ProviderMain, info.Provider)
	})
}
