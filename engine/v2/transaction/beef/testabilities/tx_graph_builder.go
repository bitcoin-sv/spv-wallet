package testabilities

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/stretchr/testify/require"
)

const testXPriv = "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3"

const (
	defaultBlockHeight = 1000
	defaultSatoshis    = 100
)

// TxDescriptor contains a transaction and its associated hexadecimal data.
type TxDescriptor struct {
	tx      *sdk.Transaction // The actual transaction.
	beefHex *string          // Hexadecimal representation of the transaction in BEEF format.
	rawHex  *string          // Optional raw hexadecimal representation of the transaction.
}

// TxID returns the actual transaction ID.
func (t TxDescriptor) TxID() string { return t.tx.TxID().String() }

// BEEFHex returns the actual transaction HexBEEF format set during TxDescriptor intialization.
func (t TxDescriptor) BEEFHex() *string { return t.beefHex }

// RawHex returns the actual transaction Hex format set during TxDescriptor intialization.
func (t TxDescriptor) RawHex() *string { return t.rawHex }

// IsBEEF checks if the actual transaction was serialized to HexBEEF format.
func (t TxDescriptor) IsBEEF() bool { return t.beefHex != nil }

// IsRaw checks if the actual transaction was serialized to Hex format.
func (t TxDescriptor) IsRaw() bool { return t.rawHex != nil }

// GraphBuilderTxs is a map where the key is a transaction name (string)
// and the value is a TxDescriptor containing details of the transaction.
type GraphBuilderTxs map[string]TxDescriptor

// Has returns true if the specified transaction node name exists in the transaction map.
func (g GraphBuilderTxs) Has(t *testing.T, name string) bool {
	t.Helper()
	require.NotEmpty(t, name, "Name of the raw transaction should not be empty")

	_, ok := g[name]
	return ok
}

// AddRawTx adds a raw transaction to the GraphBuilderTxs map.
// The name is used as the key, and the transaction is stored with its hexadecimal representation.
func (g GraphBuilderTxs) AddRawTx(t *testing.T, name string, tx *sdk.Transaction) {
	t.Helper()
	require.NotNil(t, tx, "Added transaction should not be nil")

	hex := tx.Hex()
	g[name] = TxDescriptor{tx: tx, rawHex: &hex}
}

// AddMinedTx adds a mined transaction to the GraphBuilderTxs map.
// The name is used as the key, and the transaction is stored with its BEEF hexadecimal representation.
func (g GraphBuilderTxs) AddMinedTx(t *testing.T, name string, tx *sdk.Transaction) {
	t.Helper()
	require.NotNil(t, tx, "Added transaction should not be nil")

	hex, err := tx.BEEFHex()
	require.NoError(t, err, "Failed to generate BEEF hexadecimal for transaction")

	g[name] = TxDescriptor{tx: tx, beefHex: &hex}
}

// TxGraphBuilder is a utility for constructing and managing Bitcoin transactions.
// It supports creating mined and raw transactions, verifying scripts, and managing state.
type TxGraphBuilder struct {
	t                    *testing.T      // Test context for handling errors.
	transactions         GraphBuilderTxs // Collection of transactions managed by the builder.
	p2pKHLockingScript   *script.Script  // Locking script for P2PKH transactions.
	p2pKHUnlockingScript *p2pkh.P2PKH    // Unlocking script for P2PKH transactions.
	hexGen               HexGen          // Generator for unique hexadecimal values.
	blockHeight          BlockHeight     // Current block height.
	satoshis             uint64          // Default amount of satoshis for mined transactions.
}

// Transactions returns the map of transactions managed by the builder.
func (t *TxGraphBuilder) Transactions() GraphBuilderTxs { return t.transactions }

// CreateMinedTx creates a mined transaction with the specified number of outputs.
// The transaction is initialized with default block height, satoshi values, and scripts.
func (t *TxGraphBuilder) CreateMinedTx(name string, outputs uint32) *sdk.Transaction {
	t.t.Helper()

	if t.transactions.Has(t.t, name) {
		require.FailNow(t.t, "Transaction node name exists in the transaction map")
	}

	tx := sdk.NewTransaction()
	for vout := uint32(0); vout < outputs; vout++ {
		utxo, err := sdk.NewUTXO(t.hexGen.Val(), vout, t.p2pKHLockingScript.String(), t.satoshis)
		require.NoError(t.t, err, "Failed to initialize new UTXO")
		require.NoError(t.t, tx.AddInputsFromUTXOs(utxo), "Failed to add inputs from UTXO")

		utxo.UnlockingScriptTemplate = t.p2pKHUnlockingScript
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: t.p2pKHLockingScript})
		t.hexGen.Inc()
	}

	merklePath := sdk.NewMerklePath(t.blockHeight.Val(), [][]*sdk.PathElement{{{Hash: tx.TxID(), Offset: 0}}})
	require.NoError(t.t, tx.AddMerkleProof(merklePath), "Failed to add Merkle Proof")
	require.NoError(t.t, tx.Sign(), "Failed to sign transaction")

	t.blockHeight.Inc()
	t.transactions.AddMinedTx(t.t, name, tx)
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
func (t *TxGraphBuilder) CreateRawTx(name string, parents ...ParentTx) *sdk.Transaction {
	t.t.Helper()

	if t.transactions.Has(t.t, name) {
		require.FailNow(t.t, "Transaction node name exists in the transaction map")
	}

	tx := sdk.NewTransaction()
	for _, parent := range parents {
		tx.AddInputFromTx(parent.Tx, parent.Vout, t.p2pKHUnlockingScript)
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: t.p2pKHLockingScript})
	}
	require.NoError(t.t, tx.Sign(), "Failed to sign transaction")

	t.transactions.AddRawTx(t.t, name, tx)
	return tx
}

// EnsureGraphIsValid verifies that the unlocking scripts for all inputs in the transaction
// are valid and correspond to the locking scripts of the referenced UTXOs.
// The method ensures that the constructed graph is in a valid state and can be used
// as part of a test scenario to be created.
func (t *TxGraphBuilder) EnsureGraphIsValid() {
	t.t.Helper()

	for node, descriptor := range t.transactions {
		if descriptor.IsBEEF() {
			continue
		}

		for i, input := range descriptor.tx.Inputs {
			verified, err := spv.VerifyScripts(input.SourceTransaction)
			require.NoErrorf(t.t, err, "Tx node %s failed to verify input at index %d", node, i)
			require.Truef(t.t, verified, "Tx node %s script verification failed for input at index %d", node, i)
		}
	}
}

// ToTxQueryResultSlice converts graph builder transactions to TxQueryResultSlice.
func (t *TxGraphBuilder) ToTxQueryResultSlice() beef.TxQueryResultSlice {
	var slice beef.TxQueryResultSlice
	for _, desc := range t.transactions {
		sourceTXID := desc.TxID()
		if desc.IsBEEF() {
			slice = append(slice, &beef.TxQueryResult{SourceTXID: sourceTXID, BeefHex: desc.BEEFHex()})
			continue
		}
		slice = append(slice, &beef.TxQueryResult{SourceTXID: sourceTXID, RawHex: desc.RawHex()})
	}
	return slice
}

// NewTxGraphBuilder initializes a new transaction graph builder with the provided testing context and scripts.
func NewTxGraphBuilder(t *testing.T) *TxGraphBuilder {
	scripts := NewTxScripts(t, testXPriv)
	return &TxGraphBuilder{
		t:                    t,
		transactions:         make(GraphBuilderTxs),
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
