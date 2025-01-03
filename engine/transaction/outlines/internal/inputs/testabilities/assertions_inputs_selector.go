package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InputsSelectorAssertions interface {
	WithoutError(error) SuccessfullySelectedInputsAssertions
}

type SuccessfullySelectedInputsAssertions interface {
	SelectedInputs(inputs []*database.UserUtxos) SelectedInputsAssertions
}

type SelectedInputsAssertions interface {
	AreEmpty()
	ComparingTo(inputs []*database.UserUtxos) ComparingSelectedInputsAssertions
}

type ComparingSelectedInputsAssertions interface {
	AreEntries(expectedIndexes []int)
}

type assertion struct {
	t               testing.TB
	require         *require.Assertions
	assert          *assert.Assertions
	actual          []*database.UserUtxos
	comparingSource []*database.UserUtxos
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

func (a assertion) SelectedInputs(inputs []*database.UserUtxos) SelectedInputsAssertions {
	a.t.Helper()
	a.actual = inputs
	return a
}

func (a assertion) AreEmpty() {
	a.t.Helper()
	a.require.Empty(a.actual)
}

func (a assertion) ComparingTo(inputs []*database.UserUtxos) ComparingSelectedInputsAssertions {
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
		a.assert.Equal(expectedUTXO.UserID, selectedUTXO.UserID)
		a.assert.Equal(expectedUTXO.TxID, selectedUTXO.TxID)
		a.assert.EqualValues(expectedUTXO.Vout, selectedUTXO.Vout)
	}
}
