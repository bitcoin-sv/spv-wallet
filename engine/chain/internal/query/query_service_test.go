package query

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)))
	deepSuffix, _ := uuid.NewUUID()
	deploymentID := "spv-wallet-" + deepSuffix.String()
	arcURL := "https://api.taal.com/arc"
	token := "mainnet_3af382fadbc448b15cc4133242ac2621"

	newService := func() *Service {
		return NewQueryService(logger, resty.New(), arcURL, token, deploymentID)
	}

	t.Run("Query for MINED transaction", func(t *testing.T) {
		service := newService()
		txInfo, err := service.Query(ctx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.NoError(t, err)
		require.NotNil(t, txInfo)

		require.Equal(t, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06", txInfo.TxID)
		require.Equal(t, "MINED", string(txInfo.TXStatus))
	})

	t.Run("Query for unknown transaction", func(t *testing.T) {
		service := newService()
		txInfo, err := service.Query(ctx, "aaaa1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.NoError(t, err)
		require.Nil(t, txInfo)
	})

	t.Run("Query for invalid transaction", func(t *testing.T) {
		service := newService()
		txInfo, err := service.Query(ctx, "invalid")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrInvalidTransactionID)
		require.Nil(t, txInfo)
	})

	t.Run("Query with wrong token", func(t *testing.T) {
		service := newService()
		service.token = "wrong-token"
		_, err := service.Query(ctx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnauthorized)
	})

	t.Run("Query 404 endpoint but reachable", func(t *testing.T) {
		service := newService()
		service.url = "https://api.taal.com/arc-wrong"
		_, err := service.Query(ctx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
	})

	t.Run("Query 404 endpoint with wrong arcURL", func(t *testing.T) {
		service := newService()
		service.url = "wrong-arcURL"
		_, err := service.Query(ctx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
	})

	t.Run("Query interrupted by ctx timeout", func(t *testing.T) {
		service := newService()
		tCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()
		_, err := service.Query(tCtx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
	})

	t.Run("Query interrupted by resty timeout", func(t *testing.T) {
		service := newService()
		service.httpClient.SetTimeout(1 * time.Millisecond)
		_, err := service.Query(ctx, "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrARCUnreachable)
	})
}
