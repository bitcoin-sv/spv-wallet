package sql_test

import (
	"context"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo/internal/sql/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

func TestInputsSelector(t *testing.T) {

	t.Run("return empty list when db is empty", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		selector := given.NewInputSelector()

		// when:
		utxos, err := selector.Select(context.Background(), sdk.NewTransaction(), fixtures.Sender.ID())

		// then:
		then.WithoutError(err).SelectedInputs(utxos).AreEmpty()
	})

	singleSelectTests := map[string]struct {
		selectBy             selectBy
		expectToSelectInputs []int
	}{
		"select empty list when user has not enough funds": {
			selectBy: selectBy{
				satoshis: 1_000_000,
			},
		},
		"select inputs that covers outputs and fee without change": {
			selectBy: selectBy{
				satoshis: 9,
			},
			expectToSelectInputs: []int{0},
		},
		"select inputs that covers outputs and fee with change": {
			selectBy: selectBy{
				satoshis: 15,
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select more inputs with change when satoshis are equal to single utxo": {
			selectBy: selectBy{
				satoshis: 10,
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select inputs that covers outputs and fee for data requiring more fee": {
			selectBy: selectBy{
				satoshis:            9,
				txSizeWithoutInputs: testabilities.MaxSizeWithoutFeeForSingleInput + 1,
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select inputs when size is equal to fee unit bytes": {
			selectBy: selectBy{
				satoshis:            9,
				txSizeWithoutInputs: testabilities.MaxSizeWithoutFeeForSingleInput,
			},
			expectToSelectInputs: []int{0},
		},
	}
	for name, test := range singleSelectTests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and: having some utxo in database
			ownedInputs := []*database.UserUTXO{
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
			}

			// and:
			bsvTransaction := given.Transaction().ForSatoshisAndSize(&test.selectBy)

			// and:
			selector := given.NewInputSelector()

			// when:
			utxos, err := selector.Select(context.Background(), bsvTransaction, fixtures.Sender.ID())

			// then:
			then.WithoutError(err).SelectedInputs(utxos).
				ComparingTo(ownedInputs).AreEntries(test.expectToSelectInputs)

		})
	}

	twiceSelectTests := map[string]struct {
		selectBy             selectBy
		expectToSelectInputs []int
	}{
		"select different inputs for second call": {
			selectBy: selectBy{
				satoshis: 15,
			},
			expectToSelectInputs: []int{2, 3},
		},
		"select already touched inputs if the amount of not touched won't fulfill required amount": {
			selectBy: selectBy{
				satoshis: 25,
			},
			expectToSelectInputs: []int{0, 1, 3},
		},
	}
	for name, test := range twiceSelectTests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and: having some utxo in database
			ownedInputs := []*database.UserUTXO{
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
				given.DB().HasUTXO().OwnedBySender().P2PKH().WithSatoshis(10).Stored(),
			}

			// and:
			bsvTransaction := given.Transaction().ForSatoshisAndSize(&test.selectBy)

			// and:
			selector := given.NewInputSelector()

			// when:
			_, err := selector.Select(context.Background(), bsvTransaction, fixtures.Sender.ID())

			// then:
			require.NoError(t, err)

			// when:
			utxos, err := selector.Select(context.Background(), bsvTransaction, fixtures.Sender.ID())

			// then:
			then.WithoutError(err).SelectedInputs(utxos).
				ComparingTo(ownedInputs).AreEntries(test.expectToSelectInputs)

		})
	}
}

type selectBy struct {
	satoshis            bsv.Satoshis
	txSizeWithoutInputs int
}

func (s *selectBy) Satoshis() bsv.Satoshis {
	return s.satoshis
}

func (s *selectBy) Size() int {
	if s.txSizeWithoutInputs == 0 {
		return testabilities.SizeOfTransactionWithOnlyP2PKHOutput
	}
	return s.txSizeWithoutInputs
}
