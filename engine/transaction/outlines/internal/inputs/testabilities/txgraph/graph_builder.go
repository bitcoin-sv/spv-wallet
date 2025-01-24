package txgraph

import (
	"strings"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

type TxNode string

func (t TxNode) IsMinedTx() bool { return strings.Contains(string(t), "Mined") }

func (t TxNode) IsRawTx() bool { return !t.IsMinedTx() }

type TxNodes map[TxNode]*sdk.Transaction

func (t TxNodes) Add(node TxNode, tx *sdk.Transaction) { t[node] = tx }

func (t TxNodes) Has(node TxNode) bool {
	_, ok := t[node]
	return ok
}

func (t TxNodes) Val(node TxNode) *sdk.Transaction {
	v, _ := t[node]
	return v
}

type GraphBuilder struct {
	T               *testing.T
	RawTxBuilder    *RawTxBuilder
	MinedxTxBuilder *MinedxTxBuilder
	TxNodes         TxNodes
}

type TxGraph map[TxNode][]TxNode

func (t TxGraph) Has(node TxNode) bool {
	_, ok := t[node]
	return ok
}

func (g *GraphBuilder) Build(graph TxGraph) *sdk.Transaction {
	g.T.Helper()

	if !graph.Has("tx0") {
		require.FailNow(g.T, "TxGraph must contain Tx0 as starting point")
	}

	var buildTx func(node TxNode) *sdk.Transaction
	buildTx = func(node TxNode) *sdk.Transaction {
		if g.TxNodes.Has(node) {
			return g.TxNodes.Val(node)
		}

		var tx *sdk.Transaction
		if node.IsMinedTx() {
			tx = g.MinedxTxBuilder.MakeTx(1) // inputs slice size
			g.TxNodes.Add(node, tx)
			return tx
		}

		var ascendants []AscendantTx
		for _, ascendant := range graph[node] {
			ascendantTx := buildTx(ascendant)
			ascendants = append(ascendants, AscendantTx{
				Tx:   ascendantTx,
				Vout: 0,
			})
		}
		tx = g.RawTxBuilder.MakeTx(ascendants...)

		g.TxNodes.Add(node, tx)
		return tx
	}

	return buildTx("tx0") // starting point
}

func NewGraphBuilder(t *testing.T, scripts *TxScriptsBuilder) *GraphBuilder {
	return &GraphBuilder{
		T: t,
		RawTxBuilder: &RawTxBuilder{
			T:                    t,
			P2PKHLockingScript:   scripts.P2PKHLockingScript(),
			P2PKHUnlockingScript: scripts.P2PKHUnlockingScriptTemplate(),
		},
		MinedxTxBuilder: &MinedxTxBuilder{
			T:                    t,
			P2PKHLockingScript:   scripts.P2PKHLockingScript(),
			P2PKHUnlockingScript: scripts.P2PKHUnlockingScriptTemplate(),
			Block:                1000,
			Satoshis:             10,
		},
		TxNodes: make(map[TxNode]*sdk.Transaction),
	}
}
