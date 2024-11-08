package fixtures

import (
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

const (
	// XPrivStr is the xpriv key which can be safely used for testing purposes
	XPrivStr = "xprv9s21ZrQH143K2mYcVbQKRnVp8r8cVUvNVkfk7wvMd9rPxf86SoW4dJTH7Zduta6tuYS1XEKTb8hxK4zuz2fSVW9iikoRcWor8d2mGjtmke5"
)

// Fixtures for testing purposes
var (
	XPriv, _                        = bip32.GenerateHDKeyFromString(XPrivStr)
	Priv, _                         = bip32.GetPrivateKeyFromHDKey(XPriv)
	Pub                             = Priv.PubKey()
	Address, _                      = script.NewAddressFromPublicKey(Pub, true)
	P2PKHLockingScript, _           = p2pkh.Lock(Address)
	P2PKHUnlockingScriptTemplate, _ = p2pkh.Unlock(Priv, ptr(sighash.AllForkID))
)

// grandparentTXIDs are used to indicate prevTXID for parentTXs(source transactions)
// [grandparentTX] -> [parentTX] -> [tx]
// tx is the actual transaction that is being created
// parentTX contains merkle proof
// parentOfSource txs are just placeholders (IDs only)
var grandparentTXIDs = []string{
	"56acf2ec5ae4906e8814369849f60d07ba339671e57f52200ab473e183463f2e",
	"6b0fc7403ffa214357f0326224903e612acf5c3fc5b88dfc175a2be81e343609",
	"6bed3bef2bb8b41289bca3e7f92fa5e7714decb404590ccbbc6ff7dcabf0c725",
	"e4fb03a7eae2f3766f76d52719633a65c7882c7734ae7afe603e97d193f42c0e",
	"53e784cb876751e114c9e4c1921240f184b6dff8d167715227e3708d0f7bb26d",
	"27a53423aa3e5d5c46bf30be53a9998dd247daf758847f244f82d430be71de6e",
}

// GivenTXSpec is a builder for creating MOCK! transactions
type GivenTXSpec interface {
	WithoutSigning() GivenTXSpec
	WithInput(satoshis uint64) GivenTXSpec
	WithSingleSourceInputs(satoshis ...uint64) GivenTXSpec
	WithOPReturn(dataStr string) GivenTXSpec
	WithOutputScript(parts ...ScriptPart) GivenTXSpec
	WithP2PKHOutput(satoshis uint64) GivenTXSpec

	TX() *trx.Transaction
	InputUTXO(inputID int) bsv.Outpoint
	InputSourceTX(inputID int) *trx.Transaction
	ID() string
	BEEF() string
	RawTX() string
	EF() string
}

type txSpec struct {
	utxos          []*trx.UTXO
	outputs        []*trx.TransactionOutput
	t              testing.TB
	disableSigning bool

	grandparentTXIndex int
	sourceTransactions map[string]*trx.Transaction
}

// GivenTX creates a new GivenTXSpec for building a MOCK! transaction
func GivenTX(t testing.TB) GivenTXSpec {
	return &txSpec{
		t:                  t,
		grandparentTXIndex: 0,
		sourceTransactions: make(map[string]*trx.Transaction),
	}
}

// WithoutSigning disables signing of the transaction (default is false)
func (spec *txSpec) WithoutSigning() GivenTXSpec {
	spec.disableSigning = true
	return spec
}

// WithInput adds an input to the transaction with the specified satoshis
// it automatically creates a parent tx (sourceTX) with P2PKH UTXO with provided satoshis
func (spec *txSpec) WithInput(satoshis uint64) GivenTXSpec {
	return spec.WithSingleSourceInputs(satoshis)
}

// WithSingleSourceInputs adds inputs to the transaction with the specified satoshis
// All the inputs will be sourced from a single parent transaction
func (spec *txSpec) WithSingleSourceInputs(satoshis ...uint64) GivenTXSpec {
	sourceTX := spec.makeParentTX(satoshis...)
	for i, s := range satoshis {
		i32, _ := conv.IntToUint32(i)
		utxo, err := trx.NewUTXO(sourceTX.TxID().String(), i32, P2PKHLockingScript.String(), s)
		require.NoErrorf(spec.t, err, "creating utxo for input: %d", i)

		utxo.UnlockingScriptTemplate = P2PKHUnlockingScriptTemplate

		spec.utxos = append(spec.utxos, utxo)
	}
	spec.sourceTransactions[sourceTX.TxID().String()] = sourceTX

	return spec
}

// WithOPReturn adds an OP_RETURN output to the transaction with the specified data
func (spec *txSpec) WithOPReturn(dataStr string) GivenTXSpec {
	data := []byte(dataStr)
	o, err := trx.CreateOpReturnOutput([][]byte{data})
	require.NoError(spec.t, err, "creating op return output")

	spec.outputs = append(spec.outputs, o)

	return spec
}

// WithP2PKHOutput adds a P2PKH output to the transaction with the specified satoshis
func (spec *txSpec) WithP2PKHOutput(satoshis uint64) GivenTXSpec {
	spec.outputs = append(spec.outputs, &trx.TransactionOutput{
		Satoshis:      satoshis,
		LockingScript: P2PKHLockingScript,
	})
	return spec
}

// ScriptPart is an interface for building script parts
type ScriptPart interface {
	Append(s *script.Script)
}

// OpCode is an alias for byte to represent an opcode and implements ScriptPart for script building
type OpCode byte

// Append appends the opcode to the script
func (op OpCode) Append(s *script.Script) {
	_ = s.AppendOpcodes(script.OpRETURN)
}

// PushData is an alias for []byte to represent push data and implements ScriptPart for script building
type PushData []byte

// Append appends the push data to the script
func (data PushData) Append(s *script.Script) {
	_ = s.AppendPushData(data)
}

// WithOutputScript adds an output to the transaction with the specified script parts
func (spec *txSpec) WithOutputScript(parts ...ScriptPart) GivenTXSpec {
	s := &script.Script{}
	for _, part := range parts {
		part.Append(s)
	}
	spec.outputs = append(spec.outputs, &trx.TransactionOutput{LockingScript: s})
	return spec
}

// TX creates a go-sdk transaction from the given spec
func (spec *txSpec) TX() *trx.Transaction {
	tx := trx.NewTransaction()
	err := tx.AddInputsFromUTXOs(spec.utxos...)
	require.NoError(spec.t, err, "adding inputs to tx")

	for _, output := range spec.outputs {
		tx.AddOutput(output)
	}

	for _, input := range tx.Inputs {
		if sourceTX := spec.sourceTransactions[input.SourceTXID.String()]; sourceTX != nil {
			input.SourceTransaction = sourceTX
		}
	}

	if !spec.disableSigning {
		err = tx.Sign()
		require.NoError(spec.t, err, "signing tx")
	}
	return tx
}

// InputUTXO returns UTXO outpoint for the input with the specified index
func (spec *txSpec) InputUTXO(inputID int) bsv.Outpoint {
	return bsv.Outpoint{
		TxID: spec.utxos[inputID].TxID.String(),
		Vout: spec.utxos[inputID].Vout,
	}
}

// InputSourceTX returns the source transaction for the input with the specified index
func (spec *txSpec) InputSourceTX(inputID int) *trx.Transaction {
	return spec.sourceTransactions[spec.utxos[inputID].TxID.String()]
}

// ID returns the transaction ID
func (spec *txSpec) ID() string {
	return spec.TX().TxID().String()
}

// BEEF returns the BEEF hex of the transaction
func (spec *txSpec) BEEF() string {
	tx := spec.TX()
	beef, err := tx.BEEFHex()
	require.NoError(spec.t, err, "getting beef hex")

	return beef
}

// RawTX returns the raw hex of the transaction
func (spec *txSpec) RawTX() string {
	tx := spec.TX()
	return tx.Hex()
}

// EF returns the EF hex of the transaction
func (spec *txSpec) EF() string {
	tx := spec.TX()
	ef, err := tx.EFHex()
	require.NoError(spec.t, err, "getting ef hex")

	return ef
}

func (spec *txSpec) makeParentTX(satoshis ...uint64) *trx.Transaction {
	tx := trx.NewTransaction()

	total := uint64(0)
	for _, s := range satoshis {
		total += s
	}
	withFee := total + 1
	utxo, err := trx.NewUTXO(spec.getNextGrandparentTXID(), 0, P2PKHLockingScript.String(), withFee)
	require.NoError(spec.t, err, "creating utxo for parent tx")

	utxo.UnlockingScriptTemplate = P2PKHUnlockingScriptTemplate

	err = tx.AddInputsFromUTXOs(utxo)
	require.NoError(spec.t, err, "adding input to parent tx")

	for _, s := range satoshis {
		tx.AddOutput(&trx.TransactionOutput{
			Satoshis:      s,
			LockingScript: P2PKHLockingScript,
		})
	}
	err = tx.Sign()
	require.NoError(spec.t, err, "signing parent tx")

	err = tx.AddMerkleProof(trx.NewMerklePath(1000, [][]*trx.PathElement{{
		&trx.PathElement{
			Hash:   tx.TxID(),
			Offset: 0,
		},
	}}))
	require.NoError(spec.t, err, "adding merkle proof to parent tx")

	return tx
}

func (spec *txSpec) getNextGrandparentTXID() string {
	id := grandparentTXIDs[spec.grandparentTXIndex]
	spec.grandparentTXIndex = (spec.grandparentTXIndex + 1) % len(grandparentTXIDs)
	return id
}

func ptr[T any](value T) *T {
	return &value
}
