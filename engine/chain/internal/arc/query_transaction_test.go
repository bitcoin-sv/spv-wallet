package arc_test

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual ARC server
*/

func TestQueryService(t *testing.T) {
	t.Run("QueryTransaction for MINED transaction", func(t *testing.T) {
		httpClient := arcMockActivate(false)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		txInfo, err := service.QueryTransaction(context.Background(), minedTxID)

		require.NoError(t, err)
		require.NotNil(t, txInfo)
		require.Equal(t, minedTxID, txInfo.TxID)
		require.Equal(t, "MINED", string(txInfo.TXStatus))
	})

	t.Run("QueryTransaction for unknown transaction", func(t *testing.T) {
		httpClient := arcMockActivate(false)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

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
			txID:      invalidTxID,
			arcToken:  arcToken,
			arcURL:    arcURL,
			expectErr: spverrors.ErrARCGenericError,
		},
		"QueryTransaction with wrong token": {
			txID:      minedTxID,
			arcToken:  "wrong-token", //if you test it on actual ARC server, this test might fail if the ARC doesn't require token
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
			httpClient := arcMockActivate(false)

			service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(tc.arcURL, tc.arcToken), chainmodels.BHSConfig{})

			txInfo, err := service.QueryTransaction(context.Background(), tc.txID)

			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectErr)
			require.Nil(t, txInfo)
		})
	}
}

func TestQueryServiceTimeouts(t *testing.T) {
	t.Run("QueryTransaction interrupted by ctx timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		txInfo, err := service.QueryTransaction(ctx, minedTxID)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})

	t.Run("QueryTransaction interrupted by resty timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		txInfo, err := service.QueryTransaction(context.Background(), minedTxID)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})
}
