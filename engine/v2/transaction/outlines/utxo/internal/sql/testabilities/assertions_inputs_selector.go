package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InputsSelectorAssertions interface {
	WithoutError(error) SuccessfullySelectedInputsAssertions
}

type SuccessfullySelectedInputsAssertions interface {
	SelectedInputs(inputs []*outlines.UTXO) SelectedInputsAssertions
}

type SelectedInputsAssertions interface {
	AreEmpty()
	ComparingTo(inputs []*database.UserUTXO) ComparingSelectedInputsAssertions
}

type ComparingSelectedInputsAssertions interface {
	AreEntries(expectedIndexes []int)
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

func (a assertion) AreEntries(expectedIndexes []int) {
	a.t.Helper()
	a.require.Len(a.actual, len(expectedIndexes))

	for i, ownedIdx := range expectedIndexes {
		selectedUTXO := a.actual[i]
		expectedUTXO := a.comparingSource[ownedIdx]
		a.assert.Equalf(expectedUTXO.TxID, selectedUTXO.TxID, "Selected different TxID at index %d", i)
		a.assert.EqualValuesf(expectedUTXO.Vout, selectedUTXO.Vout, "Selected different vout at index %d", i)
		a.assert.Equalf(bsv.CustomInstructions(expectedUTXO.CustomInstructions), selectedUTXO.CustomInstructions, "Selected different custom instructions at index %d", i)
	}
}
