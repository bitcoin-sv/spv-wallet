package txgraph_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities/txgraph"
)

func TestGraphBuilder_MinedTxGrandparentForTx1(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := builder.CreateMinedTx("tx6", 1)
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx0 := builder.CreateRawTx("tx0", txgraph.ParentTx{Tx: tx1, Vout: 0})

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxGrandparentsForTx1Tx2Tx3(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx4 := builder.CreateMinedTx("tx4", 1)
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx5 := builder.CreateMinedTx("tx5", 1)
	tx3 := builder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx6 := builder.CreateMinedTx("tx6", 1)
	tx2 := builder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx6, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxParentsForTx0(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx1 := builder.CreateMinedTx("tx1", 1)
	tx3 := builder.CreateMinedTx("tx3", 1)
	tx2 := builder.CreateMinedTx("tx2", 1)

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_CommonMinedTxGrandparentForTx1Tx3(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx5 := builder.CreateMinedTx("tx5", 2)
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx5, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 1})

	tx6 := builder.CreateMinedTx("tx6", 1)
	tx2 := builder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx6, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)

	// then:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_MinedTxGreatGrandparentForTx1Tx3(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := builder.CreateMinedTx("tx6", 1)
	tx4 := builder.CreateRawTx("tx4", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx7 := builder.CreateMinedTx("tx7", 1)
	tx5 := builder.CreateRawTx("tx5", txgraph.ParentTx{Tx: tx7, Vout: 0})
	tx3 := builder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx9 := builder.CreateMinedTx("tx9", 1)
	tx8 := builder.CreateRawTx("tx8", txgraph.ParentTx{Tx: tx9, Vout: 0})
	tx2 := builder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx8, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)

	// when:
	builder.VerifyScripts(tx0)
}

func TestGraphBuilder_CommonMinedTxGreatGrandparentForTx1Tx3(t *testing.T) {
	// given:
	builder := txgraph.NewGraphBuilder(t)

	// graph:
	tx6 := builder.CreateMinedTx("tx6", 2)
	tx4 := builder.CreateRawTx("tx4", txgraph.ParentTx{Tx: tx6, Vout: 0})
	tx1 := builder.CreateRawTx("tx1", txgraph.ParentTx{Tx: tx4, Vout: 0})

	tx5 := builder.CreateRawTx("tx5", txgraph.ParentTx{Tx: tx6, Vout: 1})
	tx3 := builder.CreateRawTx("tx3", txgraph.ParentTx{Tx: tx5, Vout: 0})

	tx7 := builder.CreateMinedTx("tx7", 1)
	tx2 := builder.CreateRawTx("tx2", txgraph.ParentTx{Tx: tx7, Vout: 0})

	tx0 := builder.CreateRawTx("tx0",
		txgraph.ParentTx{Tx: tx1, Vout: 0},
		txgraph.ParentTx{Tx: tx3, Vout: 0},
		txgraph.ParentTx{Tx: tx2, Vout: 0},
	)

	// when:
	builder.VerifyScripts(tx0)
}
