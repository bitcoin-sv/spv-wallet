package txgraph_test

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs/testabilities/txgraph"
	"github.com/stretchr/testify/require"
)

const (
	xPriv   = "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3"
	paymail = "john.doe@example.com"
)

func TestGraphBuilder_MinedGrandparentForTx1(t *testing.T) {
	// given:
	scripts := txgraph.NewTxScriptsBuilder(t, xPriv, paymail)
	builder := txgraph.NewGraphBuilder(t, scripts)

	// when:
	tx0 := builder.Build(txgraph.TxGraph{
		"tx0": []txgraph.TxNode{"tx1"}, // this throws panic due to nil ascendant
		// "tx1": []txgraph.TxNode{"tx6Mined"},
	})

	// then:
	require.NotNil(t, tx0)
	AssertTransactionInputs(t, tx0.Inputs)
}

func AssertTransactionInputs(t *testing.T, inputs []*sdk.TransactionInput) {
	for i, input := range inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoError(t, err, "input %v", i)
		require.True(t, verified, "input %v", i)
	}
}
