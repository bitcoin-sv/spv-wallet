package chainstate

import (
	"context"
	"testing"
	"time"

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
	t.Run("no tx ID", func(t *testing.T) {
		ctx := context.Background()
		c := NewTestClient(ctx, t, WithMinercraft(&minerCraftTxOnChain{}))

		_, err := c.QueryTransactionFastest(ctx, "", RequiredInMempool, 5*time.Second)
		require.Error(t, err)
	})

	t.Run("fastest query", func(t *testing.T) {
		ctx := context.Background()
		c := NewTestClient(ctx, t, WithMinercraft(&minerCraftTxOnChain{}))

		var txInfo *TransactionInfo
		txInfo, err := c.QueryTransactionFastest(ctx, onChainExample1TxID, RequiredInMempool, 5*time.Second)
		require.NoError(t, err)
		assert.NotNil(t, txInfo)
	})
}
