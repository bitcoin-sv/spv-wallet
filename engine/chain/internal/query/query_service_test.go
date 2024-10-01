package query

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryService(t *testing.T) {
	t.Run("QueryTransaction for MINED transaction", func(t *testing.T) {
		httpClient, reset := arcMockActivate()
		defer reset()

		service := NewQueryService(tester.Logger(t), httpClient, arcURL, arcToken, deploymentID())

		txInfo, err := service.QueryTransaction(context.Background(), minedTxID)

		require.NoError(t, err)
		require.NotNil(t, txInfo)
		require.Equal(t, minedTxID, txInfo.TxID)
		require.Equal(t, "MINED", string(txInfo.TXStatus))
	})

	t.Run("QueryTransaction for unknown transaction", func(t *testing.T) {
		httpClient, reset := arcMockActivate()
		defer reset()

		service := NewQueryService(tester.Logger(t), httpClient, arcURL, arcToken, deploymentID())

		txInfo, err := service.QueryTransaction(context.Background(), unknownTxID)

		require.NoError(t, err)
		require.Nil(t, txInfo)
	})
}

func TestQueryServiceErrorCases(t *testing.T) {
	errTestCases := map[string]struct {
		txID      string
		arcToken  string
		arcURL    string
		expectErr error
	}{
		"QueryTransaction for invalid transaction": {
			txID:      "invalid",
			arcToken:  arcToken,
			arcURL:    arcURL,
			expectErr: spverrors.ErrInvalidTransactionID,
		},
		"QueryTransaction with wrong token": {
			txID:      minedTxID,
			arcToken:  "wrong-token",
			arcURL:    arcURL,
			expectErr: spverrors.ErrARCUnauthorized,
		},
		"QueryTransaction 404 endpoint but reachable": {
			txID:      minedTxID,
			arcToken:  arcToken,
			arcURL:    arcURL + wrongButReachable,
			expectErr: spverrors.ErrARCUnreachable,
		},
		"QueryTransaction 404 endpoint with wrong arcURL": {
			txID:      minedTxID,
			arcToken:  arcToken,
			arcURL:    "wrong-url",
			expectErr: spverrors.ErrARCUnreachable,
		},
	}

	for name, tc := range errTestCases {
		t.Run(name, func(t *testing.T) {
			httpClient, reset := arcMockActivate()
			defer reset()

			service := NewQueryService(tester.Logger(t), httpClient, tc.arcURL, tc.arcToken, deploymentID())

			txInfo, err := service.QueryTransaction(context.Background(), tc.txID)

			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectErr)
			require.Nil(t, txInfo)
		})
	}
}

func TestQueryServiceTimeouts(t *testing.T) {
	t.Run("QueryTransaction interrupted by ctx timeout", func(t *testing.T) {
		httpClient, reset := arcMockActivate()
		defer reset()

		service := NewQueryService(tester.Logger(t), httpClient, arcURL, arcToken, deploymentID())

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		txInfo, err := service.QueryTransaction(ctx, minedTxID)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})

	t.Run("QueryTransaction interrupted by resty timeout", func(t *testing.T) {
		httpClient, reset := arcMockActivate()
		defer reset()

		service := NewQueryService(tester.Logger(t), httpClient, arcURL, arcToken, deploymentID())
		service.httpClient.SetTimeout(1 * time.Millisecond)

		txInfo, err := service.QueryTransaction(context.Background(), minedTxID)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})
}

func deploymentID() string {
	deepSuffix, _ := uuid.NewUUID()
	return "spv-wallet-" + deepSuffix.String()
}
