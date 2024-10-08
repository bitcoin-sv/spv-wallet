package bhs_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual ARC server
*/

func TestHealthcheckBHS(t *testing.T) {
	t.Run("BHS Healthcheck success", func(t *testing.T) {
		httpClient := bhsMockVerify(`{"confirmationState": "CONFIRMED"}`, false)
		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		err := service.HealthcheckBHS(context.Background())

		require.NoError(t, err)
	})

	t.Run("BHS Healthcheck reachable but invalid state", func(t *testing.T) {
		httpClient := bhsMockVerify(`{"confirmationState": "INVALID"}`, false)

		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		err := service.HealthcheckBHS(context.Background())

		require.ErrorIs(t, err, chainerrors.ErrBHSUnhealthy)
	})

	t.Run("BHS Healthcheck interrupted by ctx timeout", func(t *testing.T) {
		httpClient := bhsMockVerify("", true)
		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		err := service.HealthcheckBHS(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrBHSUnhealthy)
		require.ErrorIs(t, err, chainerrors.ErrBHSUnreachable)
	})
}
