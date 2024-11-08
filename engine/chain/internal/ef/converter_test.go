package ef_test

import (
	"context"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/ef"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func givenSingleINSingleOUTTXSpec(t *testing.T) fixtures.GivenTXSpec {
	return fixtures.GivenTX(t).WithInput(10).WithP2PKHOutput(1)
}

func givenMultipleOUTsTXSpec(t *testing.T) fixtures.GivenTXSpec {
	return givenSingleINSingleOUTTXSpec(t).WithInput(2)
}

func givenSingleSourceINsTXSpec(t *testing.T) fixtures.GivenTXSpec {
	return fixtures.GivenTX(t).WithSingleSourceInputs(1, 2).WithP2PKHOutput(1)
}

func TestConverterFromRawTx(t *testing.T) {
	givenSingleINSingleOUTTX := givenSingleINSingleOUTTXSpec(t)
	givenMultipleOUTsTX := givenMultipleOUTsTXSpec(t)
	givenSingleSourceINsTX := givenSingleSourceINsTXSpec(t)

	tests := map[string]struct {
		rawTx         string
		txGetter      *mockTransactionsGetter
		expectedEFHex string
	}{
		"Convert tx with one unsourced input": {
			rawTx: givenSingleINSingleOUTTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenSingleINSingleOUTTX.InputSourceTX(0).Hex(),
			}),
			expectedEFHex: givenSingleINSingleOUTTX.EF(),
		},
		"Convert tx with two unsourced inputs": {
			rawTx: givenMultipleOUTsTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenMultipleOUTsTX.InputSourceTX(0).Hex(),
				givenMultipleOUTsTX.InputSourceTX(1).Hex(),
			}),
			expectedEFHex: givenMultipleOUTsTX.EF(),
		},
		"Convert tx with two unsourced inputs from one source": {
			rawTx: givenSingleSourceINsTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenSingleSourceINsTX.InputSourceTX(0).Hex(),
				// NOTE: for inputID 1, the same source transaction is returned
			}),
			expectedEFHex: givenSingleSourceINsTX.EF(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx := fromHex(t, test.rawTx)

			converter := ef.NewConverter(test.txGetter)
			efHex, err := converter.Convert(context.Background(), tx)
			require.NoError(t, err)
			require.Equal(t, test.expectedEFHex, efHex)
		})
	}
}

func TestConverterAlreadyInEF(t *testing.T) {
	tests := map[string]struct {
		efHex string
	}{
		"Convert tx with one input": {
			efHex: givenSingleINSingleOUTTXSpec(t).EF(),
		},
		"Convert tx with two inputs": {
			efHex: givenMultipleOUTsTXSpec(t).EF(),
		},
		"Convert tx with two inputs from one source": {
			efHex: givenSingleSourceINsTXSpec(t).EF(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx, err := sdk.NewTransactionFromHex(test.efHex)
			require.NoError(t, err)
			converter := ef.NewConverter(newMockTransactionsGetter(t, []string{}))
			efHexRegenerated, err := converter.Convert(context.Background(), tx)
			require.NoError(t, err)
			require.Equal(t, test.efHex, efHexRegenerated)
		})
	}
}

func TestConverterErrorCases(t *testing.T) {
	givenSingleINSingleOUTTX := givenSingleINSingleOUTTXSpec(t)
	givenMultipleOUTsTX := givenMultipleOUTsTXSpec(t)

	tests := map[string]struct {
		rawTx     string
		txGetter  *mockTransactionsGetter
		expectErr error
	}{
		"No source tx provided by TransactionGetter": {
			rawTx:     givenSingleINSingleOUTTX.RawTX(),
			txGetter:  newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxSkip),
			expectErr: ef.ErrGetTransactions,
		},
		"Not every source tx provided by TransactionGetter": {
			rawTx: givenMultipleOUTsTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenMultipleOUTsTX.InputSourceTX(0).Hex(),
				// NOTE: for inputID 1, the source transaction is missing
			}).WithOnMissingBehavior(onMissingTxSkip),
			expectErr: ef.ErrGetTransactions,
		},
		"TransactionGetter error on missing transaction": {
			rawTx:     givenSingleINSingleOUTTX.RawTX(),
			txGetter:  newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxReturnError),
			expectErr: ef.ErrGetTransactions,
		},
		"Nil transaction returned by TransactionGetter": {
			rawTx:     givenSingleINSingleOUTTX.RawTX(),
			txGetter:  newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxAddNil),
			expectErr: ef.ErrGetTransactions,
		},
		"TransactionGetter returned more transactions than requested": {
			rawTx: givenSingleINSingleOUTTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenSingleINSingleOUTTX.InputSourceTX(0).Hex(),
				givenMultipleOUTsTX.InputSourceTX(1).Hex(),
			}).WithReturnAll(true),
			expectErr: ef.ErrGetTransactions,
		},
		"TransactionGetter not requested transactions but with correct length": {
			rawTx: givenSingleINSingleOUTTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenMultipleOUTsTX.InputSourceTX(1).Hex(),
			}).WithReturnAll(true),
			expectErr: ef.ErrGetTransactions,
		},
		"TransactionGetter duplicated transaction": {
			rawTx: givenSingleINSingleOUTTX.RawTX(),
			txGetter: newMockTransactionsGetter(t, []string{
				givenSingleINSingleOUTTX.InputSourceTX(0).Hex(),
				givenSingleINSingleOUTTX.InputSourceTX(0).Hex(),
			}).WithReturnAll(true),
			expectErr: ef.ErrGetTransactions,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx := fromHex(t, test.rawTx)

			converter := ef.NewConverter(test.txGetter)
			efHex, err := converter.Convert(context.Background(), tx)
			require.ErrorIs(t, err, test.expectErr)
			require.Empty(t, efHex)
		})
	}
}

func fromHex(t *testing.T, rawTx string) *sdk.Transaction {
	tx, err := sdk.NewTransactionFromHex(rawTx)
	require.NoError(t, err)
	return tx
}
