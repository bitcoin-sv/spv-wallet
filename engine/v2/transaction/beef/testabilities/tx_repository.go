package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/stretchr/testify/require"
)

type TxRepository struct {
	graph *TxGraphBuilder
}

func (t *TxRepository) QueryInputSources(ctx context.Context, sourceTXIDs ...string) (beef.TxQueryResultSlice, error) {
	return t.graph.ToTxQueryResultSlice(), nil
}

func NewTxRepository(t *testing.T, graph *TxGraphBuilder) *TxRepository {
	require.NotNil(t, graph, "Tx graph builder should not be nil")
	return &TxRepository{graph: graph}
}
