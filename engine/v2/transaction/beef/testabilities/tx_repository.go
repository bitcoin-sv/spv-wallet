package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/stretchr/testify/require"
)

// TxRepositoryOption defines a functional option for configuring TxRepository.
type TxRepositoryOption func(*TxRepository)

// WithInvalidQueryBeefHex forces the repository to return an invalid BEEF hex value in query results.
func WithInvalidQueryBeefHex() TxRepositoryOption {
	return func(tr *TxRepository) {
		tr.invalidBEEFHex = true
	}
}

// WithInvalidQueryRawHex forces the repository to return an invalid raw hex value in query results.
func WithInvalidQueryRawHex() TxRepositoryOption {
	return func(tr *TxRepository) {
		tr.invalidRawHex = true
	}
}

// TxRepository is a test utility that simulates querying transaction input sources.
type TxRepository struct {
	graph          *TxGraphBuilder // Graph structure used to build transaction query results.
	invalidBEEFHex bool            // Flag to return invalid BEEF hex in query results.
	invalidRawHex  bool            // Flag to return invalid raw hex in query results.
}

// QueryInputSources retrieves transaction input sources based on provided transaction IDs.
// It can return invalid BEEF hex or raw hex data based on configured options.
func (t *TxRepository) QueryInputSources(ctx context.Context, sourceTXIDs ...string) (beef.TxQueryResultSlice, error) {
	slice := t.graph.ToTxQueryResultSlice()
	invalidHex := "gggg" // Placeholder for invalid hex data.

	switch {
	case t.invalidBEEFHex:
		for _, q := range slice {
			q.BeefHex = &invalidHex
		}
		return slice, nil

	case t.invalidRawHex:
		for _, q := range slice {
			q.RawHex = &invalidHex
		}
		return slice, nil

	default:
		return slice, nil
	}
}

// NewTxRepository creates a new TxRepository instance for testing.
// It ensures the graph builder is not nil and applies any provided options.
func NewTxRepository(t *testing.T, graph *TxGraphBuilder, opts ...TxRepositoryOption) *TxRepository {
	t.Helper()
	require.NotNil(t, graph, "TxGraphBuilder should not be nil")

	txRepo := TxRepository{graph: graph}
	for _, opt := range opts {
		opt(&txRepo)
	}
	return &txRepo
}
