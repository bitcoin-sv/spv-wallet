package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TransactionDetailsAssertions interface {
	WithOutputValues(values ...bsv.Satoshis) TransactionDetailsAssertions
	OutputUnlockableBy(vout uint32, user fixtures.User) TransactionDetailsAssertions
}

type transactionAssertions struct {
	t           testing.TB
	tx          *sdk.Transaction
	require     *require.Assertions
	assert      *assert.Assertions
	annotations transaction.Annotations
}

func (a *transactionAssertions) WithOutputValues(values ...bsv.Satoshis) TransactionDetailsAssertions {
	a.t.Helper()
	a.require.Lenf(a.tx.Outputs, len(values), "Tx has less outputs then expected values")
	for i, v := range values {
		a.assert.Equal(v, bsv.Satoshis(a.tx.Outputs[i].Satoshis), "output value mismatch")
	}
	return a
}

func (a *transactionAssertions) OutputUnlockableBy(vout uint32, user fixtures.User) TransactionDetailsAssertions {
	a.t.Helper()
	a.assert.Less(vout, len(a.tx.Outputs), "there is no vout to unlock in transaction outputs")

	outputAnnotation, ok := a.annotations.Outputs[vout]
	if !ok {
		a.t.Fatalf("output %d has no required annotation", vout)
	}
	a.require.NotNil(outputAnnotation.CustomInstructions, "output %d has no custom instructions", vout)

	txtestability.Given(a.t).Tx().
		WithSender(user).
		WithInputFromUTXO(a.tx, vout, *outputAnnotation.CustomInstructions...).
		WithOPReturn("dummy data").
		TX() // during TX call, the transaction is signed. Should fail if the UTXO cannot be unlocked by the user.

	return a
}
