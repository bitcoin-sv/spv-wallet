package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InputAssertion interface {
	HasSourceTxID(id string) InputAssertion
	HasSourceVout(index int) InputAssertion
}

type txInputAssertion struct {
	t          testing.TB
	parent     *assertion
	assert     *assert.Assertions
	require    *require.Assertions
	input      *sdk.TransactionInput
	annotation any
	index      int
}

func (a *txInputAssertion) HasSourceTxID(id string) InputAssertion {
	a.t.Helper()
	a.assert.Equal(id, a.input.SourceTXID)
	return a
}

func (a *txInputAssertion) HasSourceVout(index int) InputAssertion {
	a.t.Helper()
	a.assert.Equal(index, a.input.SourceTxOutIndex)
	return a
}
