package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BEEFTransactionAssertion defines a fluent interface for asserting properties of a BEEF transaction.
type BEEFTransactionAssertion interface {
	// Created asserts that the BEEF hex string is not empty.
	Created(hex string) BEEFTransactionAssertion
	// WithNoError asserts that there were no errors during transaction preparation.
	WithNoError(err error) BEEFTransactionAssertion
	// WithParseableBEEFHEX asserts that the BEEF hex can be parsed into a transaction.
	WithParseableBEEFHEX() BEEFTransactionAssertion
	// WithSourceTransactions asserts that the source transactions of the inputs are valid.
	WithSourceTransactions() BEEFTransactionAssertion
	// HasError asserts that there was an error during transaction preparation.
	HasError(err, target error) BEEFTransactionAssertion
	// IsEmpty asserts that the created BEEF hex string during transaction preparation was empty.
	IsEmpty(beefHex string) BEEFTransactionAssertion
}

// Then initializes an assertion instance for testing.
func Then(t testing.TB) BEEFTransactionAssertion {
	return &assertion{t: t, require: require.New(t), assert: assert.New(t)}
}

// assertion provides the implementation of BEEFTransactionAssertion.
type assertion struct {
	t       testing.TB
	assert  *assert.Assertions
	require *require.Assertions
	beefHex string
	tx      *sdk.Transaction
}

// Created verifies that the given BEEF hex string is not empty.
func (a *assertion) Created(beefHex string) BEEFTransactionAssertion {
	a.t.Helper()
	a.require.NotEmpty(beefHex, "PrepareBEEF should not return an empty hex string")
	a.beefHex = beefHex
	return a
}

// IsEmpty verifies that the given BEEF hex string is empty.
func (a *assertion) IsEmpty(beefHex string) BEEFTransactionAssertion {
	a.t.Helper()
	a.require.Empty(beefHex, "PrepareBEEF should return an empty hex string")
	return a
}

// WithNoError verifies that no error was returned during the transaction preparation.
func (a *assertion) WithNoError(err error) BEEFTransactionAssertion {
	a.t.Helper()
	a.require.Nil(err, "PrepareBEEF has finished with error")
	return a
}

// HasError checks whether an error was returned while preparing a transaction.
func (a *assertion) HasError(err, target error) BEEFTransactionAssertion {
	a.t.Helper()
	a.require.ErrorIs(err, target, "PrepareBEEF should contain target error in the error chain")
	return a
}

// WithParseableBEEFHEX verifies that the stored BEEF hex can be parsed into a valid transaction.
func (a *assertion) WithParseableBEEFHEX() BEEFTransactionAssertion {
	a.t.Helper()
	tx, err := sdk.NewTransactionFromBEEFHex(a.beefHex)
	a.require.Nil(err, "Failed to create BEEF transaction from the given hex")
	a.require.NotNil(tx, "Build transaction from BEEF hex should not be nil")
	a.assert.NotZero(tx.Version, "tx version is 0 which is not acceptable by nodes")
	a.tx = tx
	return a
}

// WithSourceTransactions verifies that all inputs have valid source transactions using SPV verification.
func (a *assertion) WithSourceTransactions() BEEFTransactionAssertion {
	a.t.Helper()
	a.require.NotNil(a.tx, "Transaction must not be nil")
	a.require.NotNil(a.tx.Inputs, "Transaction inputs must not be nil")
	for _, input := range a.tx.Inputs {
		ok, err := spv.VerifyScripts(input.SourceTransaction)
		a.require.True(ok, "SPV script verification failed")
		a.require.NoError(err, "Unexpected error during SPV script verification")
	}
	return a
}
