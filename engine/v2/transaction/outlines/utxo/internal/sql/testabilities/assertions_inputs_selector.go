package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InputsSelectorAssertions interface {
	WithoutError(error) SuccessfullySelectedInputsAssertions
}

type SuccessfullySelectedInputsAssertions interface {
	SelectedInputs(inputs []*outlines.UTXO) SelectedInputsAssertions
	Change(change bsv.Satoshis) ChangeAssertions
}

type SelectedInputsAssertions interface {
	AreEmpty()
	ComparingTo(inputs []*database.UserUTXO) ComparingSelectedInputsAssertions
}

type ComparingSelectedInputsAssertions interface {
	AreEntries(expectedIndexes []int) ComparingSelectedInputsAssertions
}

type ChangeAssertions interface {
	EqualsTo(change uint)
}

type assertion struct {
	t               testing.TB
	require         *require.Assertions
	assert          *assert.Assertions
	actual          []*outlines.UTXO
	comparingSource []*database.UserUTXO
}

func newAssertions(t testing.TB) InputsSelectorAssertions {
	return assertion{
		t:       t,
		require: require.New(t),
		assert:  assert.New(t),
	}
}

func (a assertion) WithoutError(err error) SuccessfullySelectedInputsAssertions {
	a.t.Helper()
	a.require.NoError(err)
	return a
}

func (a assertion) SelectedInputs(inputs []*outlines.UTXO) SelectedInputsAssertions {
	a.t.Helper()
	a.actual = inputs
	return a
}

func (a assertion) AreEmpty() {
	a.t.Helper()
	a.require.Empty(a.actual)
}

func (a assertion) ComparingTo(inputs []*database.UserUTXO) ComparingSelectedInputsAssertions {
	a.t.Helper()
	a.comparingSource = inputs
	return a
}

func (a assertion) AreEntries(expectedIndexes []int) ComparingSelectedInputsAssertions {
	a.t.Helper()

	actualOutpoints := make(map[string]*outlines.UTXO)
	for _, utxo := range a.actual {
		actualOutpoints[bsv.Outpoint{TxID: utxo.TxID, Vout: utxo.Vout}.String()] = utxo
	}

	for i, ownedIdx := range expectedIndexes {
		expectedUTXO := a.comparingSource[ownedIdx]
		expectedOutpoint := bsv.Outpoint{TxID: expectedUTXO.TxID, Vout: expectedUTXO.Vout}
		selectedUTXO, ok := actualOutpoints[expectedOutpoint.String()]
		if !ok {
			a.assert.Failf("Didn't find expected UTXO", "Expected outpoint: %s", expectedOutpoint.String())
		} else {
			a.assert.Equalf(bsv.CustomInstructions(expectedUTXO.CustomInstructions), selectedUTXO.CustomInstructions, "Selected different custom instructions at index %d", i)
			delete(actualOutpoints, expectedOutpoint.String())
		}
	}

	a.assert.Empty(actualOutpoints, "Selected not expected UTXOs")

	if a.t.Failed() {
		expectedSelected := lo.Map(expectedIndexes, func(item int, index int) outlines.UTXO {
			return outlines.UTXO{
				TxID:               a.comparingSource[item].TxID,
				Vout:               a.comparingSource[item].Vout,
				CustomInstructions: bsv.CustomInstructions(a.comparingSource[item].CustomInstructions),
			}
		})

		selectedValues := lo.Map(a.actual, func(item *outlines.UTXO, index int) outlines.UTXO {
			return *item
		})

		a.require.EqualValuesf(expectedSelected, selectedValues, "Selected different UTXOs")
	}
	return a
}

type changeAssertion struct {
	t               testing.TB
	assert          *assert.Assertions
	comparingChange uint
}

func (a assertion) Change(change bsv.Satoshis) ChangeAssertions {
	return changeAssertion{
		t:               a.t,
		assert:          a.assert,
		comparingChange: uint(change),
	}
}

func (a changeAssertion) EqualsTo(change uint) {
	a.t.Helper()
	a.assert.EqualValues(change, a.comparingChange)
}
