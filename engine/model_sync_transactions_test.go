package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSyncTransaction_GetModelName will test the method GetModelName()
func TestSyncTransaction_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("valid name", func(t *testing.T) {
		syncTx := newSyncTransaction(testTxID, &SyncConfig{SyncOnChain: true, Broadcast: true}, New())
		require.NotNil(t, syncTx)
		assert.Equal(t, ModelSyncTransaction.String(), syncTx.GetModelName())
	})

	t.Run("missing config", func(t *testing.T) {
		syncTx := newSyncTransaction(testTxID, nil, New())
		require.Nil(t, syncTx)
	})
}
