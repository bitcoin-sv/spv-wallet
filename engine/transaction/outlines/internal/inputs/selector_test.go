package inputs_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type selectBy struct {
	satoshis            bsv.Satoshis
	txSizeWithoutInputs uint64
}

func TestInputsSelector(t *testing.T) {

	t.Run("return empty list when db is empty", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		selector := given.NewInputSelector()

		// when:
		utxos, err := selector.SelectInputsForTransaction(context.Background(), fixtures.Sender.ID(), 0, 0)

		// then:
		then.WithoutError(err).SelectedInputs(utxos).AreEmpty()
	})

	singleSelectTests := map[string]struct {
		selectBy             selectBy
		expectToSelectInputs []int
	}{
		"select empty list when user has not enough funds": {
			selectBy: selectBy{
				satoshis:            1_000_000,
				txSizeWithoutInputs: 116,
			},
		},
		"select inputs that covers outputs and fee without change": {
			selectBy: selectBy{
				satoshis:            9,
				txSizeWithoutInputs: 116,
			},
			expectToSelectInputs: []int{0},
		},
		"select inputs that covers outputs and fee with change": {
			selectBy: selectBy{
				satoshis:            15,
				txSizeWithoutInputs: 116,
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select more inputs with change when satoshis are equal to single utxo": {
			selectBy: selectBy{
				satoshis:            10,
				txSizeWithoutInputs: 116,
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select inputs that covers outputs and fee for more data": {
			selectBy: selectBy{
				satoshis:            9,
				txSizeWithoutInputs: uint64(fixtures.DefaultFeeUnit.Bytes + 1 - database.EstimatedInputSizeForP2PKH),
			},
			expectToSelectInputs: []int{0, 1},
		},
		"select inputs when size is equal to fee unit bytes": {
			selectBy: selectBy{
				satoshis:            9,
				txSizeWithoutInputs: uint64(fixtures.DefaultFeeUnit.Bytes - database.EstimatedInputSizeForP2PKH),
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
			selector := given.NewInputSelector()

			// when:
			utxos, err := selector.SelectInputsForTransaction(context.Background(), fixtures.Sender.ID(), test.selectBy.satoshis, test.selectBy.txSizeWithoutInputs)

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
				satoshis:            15,
				txSizeWithoutInputs: 116,
			},
			expectToSelectInputs: []int{2, 3},
		},
		"select already touched inputs if the amount of not touched won't fulfill required amount": {
			selectBy: selectBy{
				satoshis:            25,
				txSizeWithoutInputs: 116,
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
			selector := given.NewInputSelector()

			// when:
			_, err := selector.SelectInputsForTransaction(context.Background(), fixtures.Sender.ID(), test.selectBy.satoshis, test.selectBy.txSizeWithoutInputs)

			// then:
			require.NoError(t, err)

			// when:
			utxos, err := selector.SelectInputsForTransaction(context.Background(), fixtures.Sender.ID(), test.selectBy.satoshis, test.selectBy.txSizeWithoutInputs)

			// then:
			then.WithoutError(err).SelectedInputs(utxos).
				ComparingTo(ownedInputs).AreEntries(test.expectToSelectInputs)

		})
	}
}
