package txgraph

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

const testXPriv = "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3"

const (
	defaultBlockHeight = 1000
	defaultSatoshis    = 100
)

// TransactionDescriptor contains a transaction and its associated hexadecimal data.
type TransactionDescriptor struct {
	tx      *sdk.Transaction // The actual transaction.
	beefHex *string          // Hexadecimal representation of the transaction in BEEF format.
	rawHex  *string          // Optional raw hexadecimal representation of the transaction.
}

// TxID returns the actual transaction ID.
func (t TransactionDescriptor) TxID() string { return t.tx.TxID().String() }

// BEEFHex returns the actual transaction HexBEEF format set during TransactionDescriptor intialization.
func (t TransactionDescriptor) BEEFHex() *string { return t.beefHex }

// RawHex returns the actual transaction Hex format set during TransactionDescriptor intialization.
func (t TransactionDescriptor) RawHex() *string { return t.rawHex }

// IsBEEF checks if the actual transaction was serialized to HexBEEF format.
func (t TransactionDescriptor) IsBEEF() bool { return t.beefHex != nil }

// IsRaw checks if the actual transaction was serialized to Hex format.
func (t TransactionDescriptor) IsRaw() bool { return t.rawHex != nil }

// GraphBuilderTransactions is a map where the key is a transaction name (string)
// and the value is a TransactionDescriptor containing details of the transaction.
type GraphBuilderTransactions map[string]TransactionDescriptor

// Has returns true if the specified transaction node name exists in the transaction map.
func (g GraphBuilderTransactions) Has(name string) bool {
	_, ok := g[name]
	return ok
}

// AddRawTx adds a raw transaction to the GraphBuilderTransactions map.
// The name is used as the key, and the transaction is stored with its hexadecimal representation.
func (g GraphBuilderTransactions) AddRawTx(t *testing.T, name string, tx *sdk.Transaction) {
	t.Helper()
	require.NotEmpty(t, name, "Name of the raw transaction should not be empty")
	require.NotNil(t, tx, "Added transaction should not be nil")

	if g.Has(name) {
		require.FailNow(t, "Transaction node name exists in the transaction map")
	}

	hex := tx.Hex()
	g[name] = TransactionDescriptor{tx: tx, rawHex: &hex}
}

// AddMinedTx adds a mined transaction to the GraphBuilderTransactions map.
// The name is used as the key, and the transaction is stored with its BEEF hexadecimal representation.
func (g GraphBuilderTransactions) AddMinedTx(t *testing.T, name string, tx *sdk.Transaction) {
	t.Helper()
	require.NotEmpty(t, name, "Name of the mined transaction should not be empty")
	require.NotNil(t, tx, "Added transaction should not be nil")

	if g.Has(name) {
		require.FailNow(t, "Transaction node name exists in the transaction map")
	}

	hex, err := tx.BEEFHex()
	require.NoError(t, err, "Failed to generate BEEF hexadecimal for transaction")

	g[name] = TransactionDescriptor{tx: tx, beefHex: &hex}
}

// GraphBuilder is a utility for constructing and managing Bitcoin transactions.
// It supports creating mined and raw transactions, verifying scripts, and managing state.
type GraphBuilder struct {
	t                    *testing.T               // Test context for handling errors.
	transactions         GraphBuilderTransactions // Collection of transactions managed by the builder.
	p2pKHLockingScript   *script.Script           // Locking script for P2PKH transactions.
	p2pKHUnlockingScript *p2pkh.P2PKH             // Unlocking script for P2PKH transactions.
	hexGen               HexGen                   // Generator for unique hexadecimal values.
	blockHeight          BlockHeight              // Current block height.
	satoshis             uint64                   // Default amount of satoshis for mined transactions.
}

// Transactions returns the map of transactions managed by the builder.
func (g *GraphBuilder) Transactions() GraphBuilderTransactions {
	return g.transactions
}

// CreateMinedTx creates a mined transaction with the specified number of outputs.
// The transaction is initialized with default block height, satoshi values, and scripts.
func (g *GraphBuilder) CreateMinedTx(name string, outputs uint32) *sdk.Transaction {
	g.t.Helper()

	tx := sdk.NewTransaction()
	for vout := uint32(0); vout < outputs; vout++ {
		utxo, err := sdk.NewUTXO(g.hexGen.Val(), vout, g.p2pKHLockingScript.String(), g.satoshis)
		require.NoError(g.t, err, "Failed to initialize new UTXO")
		require.NoError(g.t, tx.AddInputsFromUTXOs(utxo), "Failed to add inputs from UTXO")

		utxo.UnlockingScriptTemplate = g.p2pKHUnlockingScript
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: g.p2pKHLockingScript})
		g.hexGen.Inc()
	}

	merklePath := sdk.NewMerklePath(g.blockHeight.Val(), [][]*sdk.PathElement{{{Hash: tx.TxID(), Offset: 0}}})
	require.NoError(g.t, tx.AddMerkleProof(merklePath), "Failed to add Merkle Proof")
	require.NoError(g.t, tx.Sign(), "Failed to sign transaction")

	g.blockHeight.Inc()
	g.transactions.AddMinedTx(g.t, name, tx)
	return tx
}

// ParentTx represents a parent transaction and a specific output (Vout)
// referenced as an input in another transaction.
type ParentTx struct {
	Vout uint32           // The output index in the parent transaction.
	Tx   *sdk.Transaction // The parent transaction.
}

// CreateRawTx creates a raw transaction that consumes the outputs of parent transactions.
// Each parent output (Vout) is added as an input to the new transaction.
func (g *GraphBuilder) CreateRawTx(name string, parents ...ParentTx) *sdk.Transaction {
	g.t.Helper()

	tx := sdk.NewTransaction()
	for _, parent := range parents {
		tx.AddInputFromTx(parent.Tx, parent.Vout, g.p2pKHUnlockingScript)
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: g.p2pKHLockingScript})
	}
	require.NoError(g.t, tx.Sign(), "Failed to sign transaction")

	g.transactions.AddRawTx(g.t, name, tx)
	return tx
}

// VerifyScripts verifies that the unlocking scripts for all inputs in the transaction
// are valid and correspond to the locking scripts of the referenced UTXOs.
func (g *GraphBuilder) VerifyScripts(tx *sdk.Transaction) {
	g.t.Helper()

	for i, input := range tx.Inputs {
		verified, err := spv.VerifyScripts(input.SourceTransaction)
		require.NoErrorf(g.t, err, "Failed to verify input at index %d", i)
		require.Truef(g.t, verified, "Script verification failed for input at index %d", i)
	}
}

// NewGraphBuilder initializes a new GraphBuilder with the provided testing context and scripts.
func NewGraphBuilder(t *testing.T) *GraphBuilder {
	scripts := NewTxScripts(t, testXPriv)
	return &GraphBuilder{
		t:                    t,
		transactions:         make(GraphBuilderTransactions),
		p2pKHLockingScript:   scripts.P2PKHLockingScript(),
		p2pKHUnlockingScript: scripts.P2PKHUnlockingScriptTemplate(),
		blockHeight:          BlockHeight{defaultBlockHeight},
		hexGen:               HexGen{defaultBlockHeight},
		satoshis:             defaultSatoshis,
	}
}

// HexGen generates unique hexadecimal values for transactions.
// The zero value for HexGen is an empty struct ready to use.
type HexGen struct {
	val int32
}

// Inc increments the hexadecimal generator to produce the next value.
func (h *HexGen) Inc() { h.val++ }

// Val returns the current hexadecimal value as a string.
func (h HexGen) Val() string {
	return fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", h.val)
}

// BlockHeight tracks the current block height for mined transactions.
// The zero value for HexGen is an empty struct ready to use.
type BlockHeight struct {
	val uint32
}

// Val returns the current block height as a uint32.
func (b BlockHeight) Val() uint32 { return b.val }

// Inc increments the block height by one.
func (b *BlockHeight) Inc() { b.val++ }
