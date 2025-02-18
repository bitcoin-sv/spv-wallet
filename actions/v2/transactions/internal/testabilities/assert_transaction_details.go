package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type TransactionDetailsAssertions interface {
	WithOutValues(values ...bsv.Satoshis) TransactionDetailsAssertions
	OutputUnlockableBy(vout int, user fixtures.User) TransactionDetailsAssertions
}

type transactionAssertions struct {
	t           testing.TB
	tx          *sdk.Transaction
	require     *require.Assertions
	annotations transaction.Annotations
}

func (a *transactionAssertions) WithOutValues(values ...bsv.Satoshis) TransactionDetailsAssertions {
	if len(values) != len(a.tx.Outputs) {
		a.t.Fatalf("expected %d outputs, got %d", len(values), len(a.tx.Outputs))
	}
	for i, v := range values {
		a.require.Equal(v, bsv.Satoshis(a.tx.Outputs[i].Satoshis), "output value mismatch")
	}
	return a
}

func (a *transactionAssertions) OutputUnlockableBy(vout int, user fixtures.User) TransactionDetailsAssertions {
	a.require.Less(vout, len(a.tx.Outputs), "output index out of range")

	outputAnnotation, ok := a.annotations.Outputs[vout]
	if !ok {
		a.t.Fatalf("output %d has no required annotation", vout)
	}
	a.require.NotNil(outputAnnotation.CustomInstructions, "output %d has no custom instructions", vout)

	fixtures.GivenTX(a.t).
		WithSender(user).
		WithInputFromUTXO(a.tx, uint32(vout), *outputAnnotation.CustomInstructions...).
		WithOPReturn("dummy data").
		TX() // during TX call, the transaction is signed. Should fail if the UTXO cannot be unlocked by the user.

	return a
}
