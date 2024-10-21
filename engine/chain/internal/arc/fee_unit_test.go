package arc_test

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual ARC server
*/

func TestFeeUnit(t *testing.T) {
	t.Run("Request for policy", func(t *testing.T) {
		httpClient := mockActivate(false)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		feeUnit, err := service.GetFeeUnit(context.Background())

		require.NoError(t, err)
		require.NotNil(t, feeUnit)
		require.Equal(t, 1000, feeUnit.Bytes)
		require.Equal(t, bsv.Satoshis(1), feeUnit.Satoshis)
	})
}

func TestFeeUnitErrorCases(t *testing.T) {
	errTestCases := map[string]struct {
		arcToken  string
		arcURL    string
		expectErr error
	}{
		"GetFeeUnit with wrong token": {
			arcToken:  "wrong-token", //if you test it on actual ARC server, this test might fail if the ARC doesn't require token
			arcURL:    arcURL,
			expectErr: chainerrors.ErrARCUnauthorized,
		},
		"GetFeeUnit 404 endpoint but reachable": {
			arcToken:  arcToken,
			arcURL:    arcURL + wrongButReachable,
			expectErr: chainerrors.ErrARCUnreachable,
		},
		"GetFeeUnit 404 endpoint with wrong arcURL": {
			arcToken:  arcToken,
			arcURL:    "wrong-url",
			expectErr: chainerrors.ErrARCUnreachable,
		},
	}

	for name, tc := range errTestCases {
		t.Run(name, func(t *testing.T) {
			httpClient := mockActivate(false)

			service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(tc.arcURL, tc.arcToken), chainmodels.BHSConfig{})

			feeUnit, err := service.GetFeeUnit(context.Background())

			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectErr)
			require.Nil(t, feeUnit)
		})
	}
}

func TestFeeUnitTimeouts(t *testing.T) {
	t.Run("GetPolicy interrupted by ctx timeout", func(t *testing.T) {
		httpClient := mockActivate(true)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		feeUnit, err := service.GetFeeUnit(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, feeUnit)
	})

	t.Run("GetPolicy interrupted by resty timeout", func(t *testing.T) {
		httpClient := mockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		feeUnit, err := service.GetFeeUnit(context.Background())

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, feeUnit)
	})
}
