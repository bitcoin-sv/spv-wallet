package junglebus_test

import (
	"context"
	chainerrors "github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/junglebus"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
	"time"
)

func TestTransactionGetter(t *testing.T) {
	tests := map[string]struct {
		requestedTxIDs []string
		expectedTxIDs  []string
	}{
		"Request for one transaction": {
			requestedTxIDs: []string{knownTx1},
			expectedTxIDs:  []string{knownTx1},
		},
		"Request for two transactions": {
			requestedTxIDs: []string{knownTx1, knownTx2},
			expectedTxIDs:  []string{knownTx1, knownTx2},
		},
		"Request for two known and one unknown transactions": {
			requestedTxIDs: []string{knownTx1, knownTx2, unknownTx},
			expectedTxIDs:  []string{knownTx1, knownTx2},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			httpClient := junglebusMockActivate(false)

			service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

			transactions, err := service.GetTransactions(context.Background(), slices.Values(test.requestedTxIDs))
			require.NoError(t, err)
			require.Len(t, transactions, len(test.expectedTxIDs))
			for i, tx := range transactions {
				require.Equal(t, test.expectedTxIDs[i], tx.TxID().String())
			}
		})
	}
}

func TestTransactionsGetterErrorCases(t *testing.T) {
	t.Run("Request for wrong txID", func(t *testing.T) {
		httpClient := junglebusMockActivate(false)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		transactions, err := service.GetTransactions(context.Background(), slices.Values([]string{wrongTxID}))
		require.ErrorIs(t, err, chainerrors.ErrJunglebusFailure)
		require.Nil(t, transactions)
	})
}

func TestTransactionGetterTimeouts(t *testing.T) {
	t.Run("TransactionGetter interrupted by ctx timeout", func(t *testing.T) {
		httpClient := junglebusMockActivate(true)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		transactions, err := service.GetTransactions(ctx, slices.Values([]string{knownTx1}))

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrInternal)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, transactions)
	})

	t.Run("TransactionGetter interrupted by resty timeout", func(t *testing.T) {
		httpClient := junglebusMockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		transactions, err := service.GetTransactions(context.Background(), slices.Values([]string{knownTx1}))

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrInternal)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, transactions)
	})
}
