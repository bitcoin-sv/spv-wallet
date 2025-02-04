package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InputAssertion interface {
	HasOutpoint(outpoint bsv.Outpoint) InputAssertion
	HasSourceTxID(id string) InputAssertion
	HasSourceVout(index int) InputAssertion
	HasCustomInstructions(instructions bsv.CustomInstructions)
}

type txInputAssertion struct {
	t          testing.TB
	parent     *assertion
	assert     *assert.Assertions
	require    *require.Assertions
	input      *sdk.TransactionInput
	annotation *transaction.InputAnnotation
	index      int
}

func (a *txInputAssertion) HasOutpoint(outpoint bsv.Outpoint) InputAssertion {
	a.t.Helper()
	return a.HasSourceTxID(outpoint.TxID).HasSourceVout(int(outpoint.Vout))
}

func (a *txInputAssertion) HasSourceTxID(id string) InputAssertion {
	a.t.Helper()
	hexID, err := chainhash.NewHashFromHex(id)
	a.require.NoError(err, "Failed to parse expected source transaction ID")
	a.assert.Equalf(hexID, a.input.SourceTXID, "Source Transaction ID mismatch")
	return a
}

func (a *txInputAssertion) HasSourceVout(index int) InputAssertion {
	a.t.Helper()
	a.assert.EqualValuesf(index, a.input.SourceTxOutIndex, "Source Transaction output index mismatch")
	return a
}

func (a *txInputAssertion) HasCustomInstructions(instructions bsv.CustomInstructions) {
	a.t.Helper()
	a.require.NotNilf(a.annotation, "Input %d has no annotation", a.index)
	a.assert.Equalf(instructions, a.annotation.CustomInstructions, "Input %d has invalid custom instructions", a.index)
}
