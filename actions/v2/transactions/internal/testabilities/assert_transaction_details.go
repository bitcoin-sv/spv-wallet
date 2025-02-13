package testabilities

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/require"
	"testing"
)

type TransactionDetailsAssertions interface {
	WithOutValues(values ...uint64) TransactionDetailsAssertions
}

type transactionAssertions struct {
	t       testing.TB
	tx      *sdk.Transaction
	require *require.Assertions
}

func (a *transactionAssertions) WithOutValues(values ...uint64) TransactionDetailsAssertions {
	if len(values) != len(a.tx.Outputs) {
		a.t.Fatalf("expected %d outputs, got %d", len(values), len(a.tx.Outputs))
	}
	for i, v := range values {
		a.require.Equal(v, a.tx.Outputs[i].Satoshis, "output value mismatch")
	}
	return a
}
