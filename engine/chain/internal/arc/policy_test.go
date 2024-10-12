package arc_test

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual ARC server
*/

func TestPolicyService(t *testing.T) {
	t.Run("Request for policy", func(t *testing.T) {
		httpClient := arcMockActivate(false)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		policy, err := service.GetPolicy(context.Background())

		require.NoError(t, err)
		require.NotNil(t, policy)
		require.Equal(t, 1000, policy.Content.MiningFee.Bytes)
		require.Equal(t, bsv.Satoshis(1), policy.Content.MiningFee.Satoshis)
		require.NotEmpty(t, policy.Timestamp)
	})
}

func TestPolicyServiceErrorCases(t *testing.T) {
	errTestCases := map[string]struct {
		arcToken  string
		arcURL    string
		expectErr error
	}{
		"GetPolicy with wrong token": {
			arcToken:  "wrong-token", //if you test it on actual ARC server, this test might fail if the ARC doesn't require token
			arcURL:    arcURL,
			expectErr: spverrors.ErrARCUnauthorized,
		},
		"GetPolicy 404 endpoint but reachable": {
			arcToken:  arcToken,
			arcURL:    arcURL + wrongButReachable,
			expectErr: spverrors.ErrARCUnreachable,
		},
		"GetPolicy 404 endpoint with wrong arcURL": {
			arcToken:  arcToken,
			arcURL:    "wrong-url",
			expectErr: spverrors.ErrARCUnreachable,
		},
	}

	for name, tc := range errTestCases {
		t.Run(name, func(t *testing.T) {
			httpClient := arcMockActivate(false)

			service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(tc.arcURL, tc.arcToken), chainmodels.BHSConfig{})

			policy, err := service.GetPolicy(context.Background())

			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectErr)
			require.Nil(t, policy)
		})
	}
}

func TestPolicyServiceTimeouts(t *testing.T) {
	t.Run("GetPolicy interrupted by ctx timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		txInfo, err := service.GetPolicy(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})

	t.Run("GetPolicy interrupted by resty timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		txInfo, err := service.GetPolicy(context.Background())

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})
}
