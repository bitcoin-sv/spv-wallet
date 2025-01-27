package txgraph_test

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs/testabilities/txgraph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphBuilder(t *testing.T) {
	const xPriv = "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3"
	const paymail = "sender@example.com"

	scripts := txgraph.NewTxScripts(t, xPriv, paymail)
	builder := txgraph.NewGraphBuilder(t, scripts)

	// Build the graph
	tx5 := builder.CreateMinedTx("tx5", 2)
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx5, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 1})

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
	)

	slice := builder.TxQueryResultSlice()
	assert.NotEmpty(t, slice)
	assertSourceTransactionScripts(t, tx0)
}

func assertSourceTransactionScripts(t *testing.T, tx *sdk.Transaction) {
	t.Helper()

	for i, input := range tx.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %v", i)
		require.True(t, verified, "input %v", i)
	}
}
