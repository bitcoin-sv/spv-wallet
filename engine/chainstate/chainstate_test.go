package chainstate

import (
	"context"
	"testing"
	"time"

	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewTestClient returns a test client
func NewTestClient(ctx context.Context, t *testing.T, opts ...ClientOps) ClientInterface {
	logger := zerolog.Nop()
	c, err := NewClient(
		ctx, append(opts, WithDebugging(), WithLogger(&logger))...,
	)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

// TestQueryTransactionFastest tests the querying for a transaction and returns the fastest response
func TestQueryTransactionFastest(t *testing.T) {
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockSuccess).
		Build()
	t.Run("no tx ID", func(t *testing.T) {
		ctx := context.Background()
		c := NewTestClient(ctx, t, WithBroadcastClient(bc))

		_, err := c.QueryTransactionFastest(ctx, "", RequiredInMempool, 5*time.Second)
		require.Error(t, err)
	})

	t.Run("fastest query", func(t *testing.T) {
		ctx := context.Background()
		c := NewTestClient(ctx, t, WithBroadcastClient(bc))

		var txInfo *TransactionInfo
		txInfo, err := c.QueryTransactionFastest(ctx, onChainExample1TxID, RequiredInMempool, 5*time.Second)
		require.NoError(t, err)
		assert.NotNil(t, txInfo)
	})
}
