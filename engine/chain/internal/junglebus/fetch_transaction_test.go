package junglebus_test

import (
	"context"
	"github.com/go-resty/resty/v2"
	"testing"
	"time"

	chainerrors "github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/junglebus"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual server
*/

func TestJunglebusFetchTransaction(t *testing.T) {
	t.Run("Request for transaction", func(t *testing.T) {
		httpClient := junglebusMockActivate(false)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		tx, err := service.FetchTransaction(context.Background(), knownTx1)

		require.NoError(t, err)
		require.NotNil(t, tx)
		require.Equal(t, bsv.Satoshis(39), bsv.Satoshis(tx.TotalOutputSatoshis()))
	})

	t.Run("Request for unknown transaction", func(t *testing.T) {
		httpClient := resty.New()

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		tx, err := service.FetchTransaction(context.Background(), unknownTx)

		require.Error(t, err)
		require.Nil(t, tx)
		require.ErrorIs(t, err, chainerrors.ErrJunglebusTxNotFound)
	})

	t.Run("Request for invalid transaction", func(t *testing.T) {
		httpClient := junglebusMockActivate(false)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		tx, err := service.FetchTransaction(context.Background(), wrongTxID)

		require.Error(t, err)
		require.Nil(t, tx)
		require.ErrorIs(t, err, chainerrors.ErrJunglebusFailure)
	})
}

func TestJunglebusFetchTransactionTimeouts(t *testing.T) {
	t.Run("FetchTransaction interrupted by ctx timeout", func(t *testing.T) {
		httpClient := junglebusMockActivate(true)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		tx, err := service.FetchTransaction(ctx, knownTx1)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrInternal)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, tx)
	})

	t.Run("FetchTransaction interrupted by resty timeout", func(t *testing.T) {
		httpClient := junglebusMockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := junglebus.NewJunglebusService(tester.Logger(t), httpClient)

		tx, err := service.FetchTransaction(context.Background(), knownTx1)

		require.Error(t, err)
		require.ErrorIs(t, err, spverrors.ErrInternal)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, tx)
	})
}
