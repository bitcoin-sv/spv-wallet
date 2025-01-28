package outlines

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities/txgraph"
	"github.com/stretchr/testify/require"
)

func TestSourceTransactionBuilder_BEEFGrandparentForTx1(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx1 := graphBuilder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx0 := graphBuilder.CreateRawTx("tx0", txgraph.ParentTx{Tx: tx1, Vout: 0})
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_BEEFGrandparentsForTx1Tx2Tx3(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx4 := graphBuilder.CreateMinedTx("tx4", 1)
	tx1 := graphBuilder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx5 := graphBuilder.CreateMinedTx("tx5", 1)
	tx3 := graphBuilder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx2 := graphBuilder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx6, Vout: 0})

	tx0 := graphBuilder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_BEEFParentsForTx0(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx1 := graphBuilder.CreateMinedTx("tx4", 1)
	tx3 := graphBuilder.CreateMinedTx("tx3", 1)
	tx2 := graphBuilder.CreateMinedTx("tx2", 1)
	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_CommonBEEFTxGrandparentForTx1Tx3(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx5 := graphBuilder.CreateMinedTx("tx5", 2)
	tx1 := graphBuilder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx5, Vout: 0})
	tx3 := graphBuilder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 1})

	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx2 := graphBuilder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx6, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_BEEFTxGreatGrandparentsForTx1Tx3(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx4 := graphBuilder.CreateRawTx("tx4", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx1 := graphBuilder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx7 := graphBuilder.CreateMinedTx("tx7", 1)
	tx5 := graphBuilder.CreateRawTx("tx5", txgraph.ParentTx{Tx: tx7, Vout: 0})
	tx3 := graphBuilder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx9 := graphBuilder.CreateMinedTx("tx9", 1)
	tx8 := graphBuilder.CreateRawTx("tx8", txgraph.ParentTx{Tx: tx9, Vout: 0})
	tx2 := graphBuilder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx8, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func TestSourceTransactionBuilder_CommonBEEFTxGreatGrandparentsForTx1Tx3(t *testing.T) {
	// given:
	graphBuilder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 2)
	tx4 := graphBuilder.CreateRawTx("tx4", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx1 := graphBuilder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx5 := graphBuilder.CreateRawTx("tx5", txgraph.ParentTx{Tx: tx6, Vout: 1})
	tx3 := graphBuilder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx7 := graphBuilder.CreateMinedTx("tx7", 1)
	tx2 := graphBuilder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx7, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)
	graphBuilder.VerifyScripts(tx0)

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	db := toTxQueryResultSlice(graphBuilder.GraphBuilderTransactions())

	// when:
	builder := SourceTransactionBuilder{Tx: subjectTx}
	err := builder.Build(db)

	// then:
	require.NoError(t, err)
}

func toTxQueryResultSlice(transactions txgraph.GraphBuilderTransactions) TxQueryResultSlice {
	var slice TxQueryResultSlice
	for _, desc := range transactions {
		sourceTXID := desc.Tx.TxID().String()
		if desc.IsBEEF() {
			slice = append(slice, &TxQueryResult{SourceTXID: sourceTXID, BeefHex: desc.BEEFHex})
			continue
		}
		slice = append(slice, &TxQueryResult{SourceTXID: sourceTXID, RawHex: desc.RawHex})
	}
	return slice
}
