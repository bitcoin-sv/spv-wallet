package engine

import (
	"fmt"
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

func Test_areParentsBroadcast(t *testing.T) {
	ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
	defer deferMe()

	opts := []ModelOps{WithClient(client)}

	tx, err := txFromHex(testTxHex, append(opts, New())...)
	require.NoError(t, err)

	txErr := tx.Save(ctx)
	require.NoError(t, txErr)

	tx2, err := txFromHex(testTx2Hex, append(opts, New())...)
	require.NoError(t, err)

	txErr = tx2.Save(ctx)
	require.NoError(t, txErr)

	tx3, err := txFromHex(testTx3Hex, append(opts, New())...)
	require.NoError(t, err)

	txErr = tx3.Save(ctx)
	require.NoError(t, txErr)

	// input of testTxID
	syncTx := newSyncTransaction("65bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10", &SyncConfig{SyncOnChain: false, Broadcast: false}, append(opts, New())...)
	syncTx.BroadcastStatus = SyncStatusComplete
	txErr = syncTx.Save(ctx)
	require.NoError(t, txErr)

	// input of testTxInID
	syncTx = newSyncTransaction("89fbccca3a5e2bfc8a161bf7f54e8cb5898e296ae8c23b620b89ed570711f931", &SyncConfig{SyncOnChain: false, Broadcast: false}, append(opts, New())...)
	txErr = syncTx.Save(ctx)
	require.NoError(t, txErr)

	type args struct {
		tx   *Transaction
		opts []ModelOps
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no parents",
			args: args{
				tx:   tx3,
				opts: opts,
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "parent not broadcast",
			args: args{
				tx:   tx2,
				opts: opts,
			},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name: "parent broadcast",
			args: args{
				tx:   tx,
				opts: opts,
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _areParentsBroadcasted(ctx, tt.args.tx, tt.args.opts...)
			if !tt.wantErr(t, err, fmt.Sprintf("areParentsBroadcast(%v, %v, %v)", ctx, tt.args.tx, tt.args.opts)) {
				return
			}
			assert.Equalf(t, tt.want, got, "areParentsBroadcast(%v, %v, %v)", ctx, tt.args.tx, tt.args.opts)
		})
	}
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
