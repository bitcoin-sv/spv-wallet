package txgraph

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs"
	"github.com/stretchr/testify/require"
)

// TransactionType defines the type of a transaction.
type TransactionType string

const (
	MinedTx TransactionType = "MinedTx"
	RawTx   TransactionType = "RawTx"
)

// Check if the transaction type is mined.
func (t TransactionType) IsMined() bool { return strings.Contains(string(t), string(MinedTx)) }

// Check if the transaction type is raw.
func (t TransactionType) IsRaw() bool { return !t.IsMined() }

// GraphBuilderTransactions is a map of transaction types to transactions.
type GraphBuilderTransactions map[TransactionType]*sdk.Transaction

// Add a raw transaction.
func (t GraphBuilderTransactions) AddRawTx(name string, tx *sdk.Transaction) {
	t[TransactionType(name)+RawTx] = tx
}

// Add a mined transaction.
func (t GraphBuilderTransactions) AddMinedTx(name string, tx *sdk.Transaction) {
	t[TransactionType(name)+MinedTx] = tx
}

// Convert transactions into a slice of TxQueryResult.
func (t GraphBuilderTransactions) ToTxQueryResultSlice(test *testing.T) inputs.TxQueryResultSlice {
	var results inputs.TxQueryResultSlice
	for txType, tx := range t {
		if txType.IsRaw() {
			hex := tx.Hex()
			results = append(results, &inputs.TxQueryResult{SourceTXID: tx.TxID().String(), RawHex: &hex})
			continue
		}

		hex, err := tx.BEEFHex()
		require.NoError(test, err, "failed to generate hex for transaction")
		results = append(results, &inputs.TxQueryResult{SourceTXID: tx.TxID().String(), BeefHex: &hex})
	}
	return results
}

type GraphBuilder struct {
	t                    *testing.T
	transactions         GraphBuilderTransactions
	p2pKHLockingScript   *script.Script
	p2pKHUnlockingScript *p2pkh.P2PKH
	hexGen               HexGen
	blockHeight          BlockHeight // Default block height including transactions.
	satoshis             uint64      // Default mined transaction satoshi value.
}

// Get all transactions as TxQueryResultSlice.
func (g *GraphBuilder) TxQueryResultSlice() inputs.TxQueryResultSlice {
	return g.transactions.ToTxQueryResultSlice(g.t)
}

// Create a mined transaction with a specific number of outputs.
func (g *GraphBuilder) CreateMinedTx(name string, outputs uint32) *sdk.Transaction {
	g.t.Helper()

	tx := sdk.NewTransaction()
	for vout := uint32(0); vout < outputs; vout++ {
		utxo, err := sdk.NewUTXO(g.hexGen.Val(), vout, g.p2pKHLockingScript.String(), g.satoshis)
		require.NoError(g.t, err, "failed to initialize new utxo")
		require.NoError(g.t, tx.AddInputsFromUTXOs(utxo), "failed to add inputs from utxo")

		utxo.UnlockingScriptTemplate = g.p2pKHUnlockingScript
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: g.p2pKHLockingScript})
		g.hexGen.Inc()
	}

	merklePath := sdk.NewMerklePath(g.blockHeight.Val(), [][]*sdk.PathElement{{{Hash: tx.TxID(), Offset: 0}}})
	require.NoError(g.t, tx.AddMerkleProof(merklePath), "failed to add Merkle Proof")
	require.NoError(g.t, tx.Sign(), "failed to sign transaction")

	g.blockHeight.Inc()
	g.transactions.AddMinedTx(name, tx)
	return tx
}

// ParentTx represents a parent transaction and the specific output (Vout)
// being referenced as an input in another transaction.
type ParentTx struct {
	Vout uint32           // The output index in the parent transaction.
	Tx   *sdk.Transaction // The parent transaction.
}

// Create a raw transaction with specific parent transactions.
func (g *GraphBuilder) CreateRawTx(name string, parents ...ParentTx) *sdk.Transaction {
	g.t.Helper()

	tx := sdk.NewTransaction()
	for _, parent := range parents {
		tx.AddInputFromTx(parent.Tx, parent.Vout, g.p2pKHUnlockingScript)
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: g.p2pKHLockingScript})
	}
	require.NoError(g.t, tx.Sign(), "failed to sign transaction")

	g.transactions.AddRawTx(name, tx)
	return tx
}

// Initialize a new GraphBuilder.
func NewGraphBuilder(t *testing.T, scripts *TxScripts) *GraphBuilder {
	return &GraphBuilder{
		t:                    t,
		transactions:         make(GraphBuilderTransactions),
		p2pKHLockingScript:   scripts.P2PKHLockingScript(),
		p2pKHUnlockingScript: scripts.P2PKHUnlockingScriptTemplate(),
		hexGen:               0,
		blockHeight:          1000,
		satoshis:             100,
	}
}

// HexGen generates unique hexadecimal values.
type HexGen int32

func (h *HexGen) Inc() { *h++ }
func (h HexGen) Val() string {
	return fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", h)
}

// BlockHeight tracks the current block height.
type BlockHeight uint32

func (b BlockHeight) Val() uint32 { return uint32(b) }
func (b *BlockHeight) Inc()       { *b++ }
