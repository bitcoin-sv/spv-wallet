package arc_test

import (
	"context"
	"testing"
	"time"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/stretchr/testify/require"
	"iter"
)

func TestBroadcastTransaction(t *testing.T) {
	tests := map[string]struct {
		hex            string
		arcCfgModifier func(cfg *chainmodels.ARCConfig)
	}{
		"Broadcast unsourced tx with txs getter provided": {
			hex: validRawHex,
			arcCfgModifier: func(cfg *chainmodels.ARCConfig) {
				cfg.TxsGetter = &mockTxsGetter{
					transactions: []*sdk.Transaction{fromHex(sourceOfValidRawHex)},
				}
			},
		},
		"Broadcast tx in EF with no txs getter": {
			hex: efOfValidRawHex,
		},
		"Broadcast unsourced tx with no txs getter - raw hex as fallback": {
			hex: fallbackRawHex,
		},
		"Broadcast two-missing-inputs unsourced tx with txs getter and junglebus": {
			hex: txWithMultipleInputs,
			arcCfgModifier: func(cfg *chainmodels.ARCConfig) {
				cfg.TxsGetter = &mockTxsGetter{
					// first missing input source is provided by this txs getter (mocking getting from database)
					transactions: []*sdk.Transaction{fromHex(sourceOneOfTxWithMultipleInputs)},
				}
				cfg.UseJunglebus = true //second missing input source is provided by junglebus (mocked)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			httpClient := arcMockActivate(false)

			tx, err := sdk.NewTransactionFromHex(test.hex)
			require.NoError(t, err)

			cfg := arcCfg(arcURL, arcToken)
			if test.arcCfgModifier != nil {
				test.arcCfgModifier(&cfg)
			}

			service := chain.NewChainService(tester.Logger(t), httpClient, cfg, chainmodels.BHSConfig{})

			txInfo, err := service.Broadcast(context.Background(), tx)
			require.NoError(t, err)
			require.Equal(t, tx.TxID().String(), txInfo.TxID)
			require.Equal(t, chainmodels.SeenOnNetwork, txInfo.TXStatus)
		})
	}
}

func TestBroadcastTransactionErrorCases(t *testing.T) {
	tests := map[string]struct {
		hex            string
		arcCfgModifier func(cfg *chainmodels.ARCConfig)
		expectErr      error
	}{
		"Double spend attempt with 'old' UTXO": {
			hex:       oldWithDoubleSpentHex,
			expectErr: chainerrors.ErrARCProblematicStatus,
		},
		"Double spend attempt with relatively 'new' UTXO": {
			hex:       newWithDoubleSpentHex,
			expectErr: chainerrors.ErrARCProblematicStatus,
		},
		"Broadcast malformed tx": {
			hex:       malformedTxHex,
			expectErr: chainerrors.ErrARCUnprocessable,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			httpClient := arcMockActivate(false)

			tx, err := sdk.NewTransactionFromHex(test.hex)
			require.NoError(t, err)

			cfg := arcCfg(arcURL, arcToken)
			if test.arcCfgModifier != nil {
				test.arcCfgModifier(&cfg)
			}

			service := chain.NewChainService(tester.Logger(t), httpClient, cfg, chainmodels.BHSConfig{})

			txInfo, err := service.Broadcast(context.Background(), tx)
			require.ErrorIs(t, err, test.expectErr)
			require.Nil(t, txInfo)
		})
	}
}

func TestBroadcastTimeouts(t *testing.T) {
	t.Run("Broadcast transaction interrupted by ctx timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)

		tx, err := sdk.NewTransactionFromHex(efOfValidRawHex)
		require.NoError(t, err)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		txInfo, err := service.Broadcast(ctx, tx)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})

	t.Run("Broadcast transaction interrupted by resty timeout", func(t *testing.T) {
		httpClient := arcMockActivate(true)
		httpClient.SetTimeout(1 * time.Millisecond)

		tx, err := sdk.NewTransactionFromHex(efOfValidRawHex)
		require.NoError(t, err)

		service := chain.NewChainService(tester.Logger(t), httpClient, arcCfg(arcURL, arcToken), chainmodels.BHSConfig{})

		txInfo, err := service.Broadcast(context.Background(), tx)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrARCUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, txInfo)
	})
}

type mockTxsGetter struct {
	transactions []*sdk.Transaction
}

func (mtg *mockTxsGetter) GetTransactions(_ context.Context, _ iter.Seq[string]) ([]*sdk.Transaction, error) {
	return mtg.transactions, nil
}

func fromHex(hex string) *sdk.Transaction {
	tx, _ := sdk.NewTransactionFromHex(hex)
	return tx
}
