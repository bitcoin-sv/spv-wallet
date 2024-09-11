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

func TestSyncTransaction_SaveHook(t *testing.T) {
	t.Parallel()

	t.Run("trim Results to last 20 messages", func(t *testing.T) {
		// Given
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		opts := []ModelOps{WithClient(client), New()}
		syncTx := newSyncTransaction(testTxID, &SyncConfig{SyncOnChain: true, Broadcast: true}, opts...)

		txErr := syncTx.Save(ctx)
		require.NoError(t, txErr)

		// When
		for i := 0; i < 40; i++ {
			syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{Action: "test", StatusMessage: "msg"})
		}
		txErr = syncTx.Save(ctx)
		require.NoError(t, txErr)

		// Then
		resultsLen := len(syncTx.Results.Results)
		require.Equal(t, 20, resultsLen)
	})
}
