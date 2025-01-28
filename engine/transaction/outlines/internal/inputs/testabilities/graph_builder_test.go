package testabilities_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs/testabilities"
)

const xPriv = "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3"
const paymail = "sender@example.com"

func TestGraphBuilder_MinedTxGrandparentForTx1(t *testing.T) {
	// given:
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph:
	tx6 := builder.CreateMinedTx("tx6", 1)
	tx1 := builder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx0 := builder.CreateRawTx("tx0", testabilities.ParentTx{Tx: tx1, Vout: 0})

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxGrandparentsForTx1Tx2Tx3(t *testing.T) {
	// given:
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph
	tx4 := builder.CreateMinedTx("tx4", 1)
	tx1 := builder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx5 := builder.CreateMinedTx("tx5", 1)
	tx3 := builder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx6 := builder.CreateMinedTx("tx6", 1)
	tx2 := builder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx6, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxParentsForTx0(t *testing.T) {
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph
	tx1 := builder.CreateMinedTx("tx1", 1)
	tx3 := builder.CreateMinedTx("tx3", 1)
	tx2 := builder.CreateMinedTx("tx2", 1)

	tx0 := builder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_CommonMinedTxGrandparentForTx1Tx3(t *testing.T) {
	// given:
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph
	tx5 := builder.CreateMinedTx("tx5", 2)
	tx1 := builder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx5, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 1})

	tx6 := builder.CreateMinedTx("tx6", 1)
	tx2 := builder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx6, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxGreatGrandparentForTx1Tx3(t *testing.T) {
	// given:
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph
	tx6 := builder.CreateMinedTx("tx6", 1)
	tx4 := builder.CreateRawTx("tx4", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx1 := builder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx7 := builder.CreateMinedTx("tx7", 1)
	tx5 := builder.CreateRawTx("tx5", testabilities.ParentTx{Tx: tx7, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx9 := builder.CreateMinedTx("tx9", 1)
	tx8 := builder.CreateRawTx("tx8", testabilities.ParentTx{Tx: tx9, Vout: 0})
	tx2 := builder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx8, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	// when:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_CommonMinedTxGreatGrandparentForTx1Tx3(t *testing.T) {
	// given:
	scripts := testabilities.NewTxScripts(t, xPriv, paymail)
	builder := testabilities.NewGraphBuilder(t, scripts)

	// Build the graph
	tx6 := builder.CreateMinedTx("tx6", 1)
	tx4 := builder.CreateRawTx("tx4", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx1 := builder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx5 := builder.CreateRawTx("tx5", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx7 := builder.CreateMinedTx("tx7", 1)
	tx2 := builder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx7, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	// when:
	builder.VerifyScripts(tx0)
}
