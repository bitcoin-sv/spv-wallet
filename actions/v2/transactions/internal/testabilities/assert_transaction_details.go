package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type TransactionDetailsAssertions interface {
	WithOutValues(values ...bsv.Satoshis) TransactionDetailsAssertions
}

type transactionAssertions struct {
	t       testing.TB
	tx      *sdk.Transaction
	require *require.Assertions
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
