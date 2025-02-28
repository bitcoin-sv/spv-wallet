package txtestability

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
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

type TransactionsFixtures interface {
	Tx() TransactionSpec
}

// TransactionSpec is a builder for creating MOCK! transactions
type TransactionSpec interface {
	WithSender(sender fixtures.User) TransactionSpec
	WithRecipient(recipient fixtures.User) TransactionSpec
	WithoutSigning() TransactionSpec
	WithInput(satoshis uint64) TransactionSpec
	WithInputFromUTXO(tx *trx.Transaction, vout uint32, customInstructions ...bsv.CustomInstruction) TransactionSpec
	WithSingleSourceInputs(satoshis ...uint64) TransactionSpec
	WithOPReturn(dataStr string) TransactionSpec
	WithOutputScriptParts(parts ...ScriptPart) TransactionSpec
	WithOutputScript(satoshis uint64, script *script.Script) TransactionSpec
	WithP2PKHOutput(satoshis uint64) TransactionSpec

	TX() *trx.Transaction
	InputUTXO(inputID int) bsv.Outpoint
	InputSourceTX(inputID int) *trx.Transaction
	ID() string
	BEEF() string
	RawTX() string
	EF() string
}

type txFixturesFactory struct {
	t                  testing.TB
	blockHeight        uint32
	grandparentTXIndex int
}

func Given(t testing.TB) TransactionsFixtures {
	return &txFixturesFactory{
		t:                  t,
		blockHeight:        1000,
		grandparentTXIndex: 0,
	}
}

func (f *txFixturesFactory) Tx() TransactionSpec {
	spec := &txSpec{
		t:                  f.t,
		sourceTransactions: make(map[string]*trx.Transaction),
		sender:             fixtures.Sender,
		recipient:          fixtures.RecipientInternal,
		fixtureFactory:     f,
	}

	return spec
}

func (f *txFixturesFactory) getNextGrandparentTXID() string {
	id := grandparentTXIDs[f.grandparentTXIndex]
	f.grandparentTXIndex = (f.grandparentTXIndex + 1) % len(grandparentTXIDs)
	return id
}

func (f *txFixturesFactory) getNextBlockHeight() uint32 {
	h := f.blockHeight
	f.blockHeight++
	return h
}

type txSpec struct {
	utxos          []*trx.UTXO
	outputs        []*trx.TransactionOutput
	t              testing.TB
	disableSigning bool

	grandparentTXIndex int
	sourceTransactions map[string]*trx.Transaction
	blockHeight        uint32
	sender             fixtures.User
	recipient          fixtures.User
	fixtureFactory     *txFixturesFactory
}

// WithSender sets the sender for the transaction (default is Sender)
func (spec *txSpec) WithSender(sender fixtures.User) TransactionSpec {
	spec.sender = sender
	return spec
}

// WithRecipient sets the recipient for the transaction (default is RecipientInternal)
func (spec *txSpec) WithRecipient(recipient fixtures.User) TransactionSpec {
	spec.recipient = recipient
	return spec
}

// WithoutSigning disables signing of the transaction (default is false)
func (spec *txSpec) WithoutSigning() TransactionSpec {
	spec.disableSigning = true
	return spec
}

// WithInput adds an input to the transaction with the specified satoshis
// it automatically creates a parent tx (sourceTX) with P2PKH UTXO with provided satoshis
func (spec *txSpec) WithInput(satoshis uint64) TransactionSpec {
	return spec.WithSingleSourceInputs(satoshis)
}

func (spec *txSpec) WithInputFromUTXO(tx *trx.Transaction, vout uint32, customInstructions ...bsv.CustomInstruction) TransactionSpec {
	output := tx.Outputs[vout]
	utxo, err := trx.NewUTXO(tx.TxID().String(), vout, output.LockingScript.String(), output.Satoshis)
	require.NoError(spec.t, err, "creating utxo for input")

	utxo.UnlockingScriptTemplate = spec.sender.P2PKHUnlockingScriptTemplate(customInstructions...)

	spec.utxos = append(spec.utxos, utxo)
	spec.sourceTransactions[tx.TxID().String()] = tx
	return spec
}

// WithSingleSourceInputs adds inputs to the transaction with the specified satoshis
// All the inputs will be sourced from a single parent transaction
func (spec *txSpec) WithSingleSourceInputs(satoshis ...uint64) TransactionSpec {
	sourceTX := spec.makeParentTX(satoshis...)
	for i, s := range satoshis {
		i32, _ := conv.IntToUint32(i)
		utxo, err := trx.NewUTXO(sourceTX.TxID().String(), i32, spec.sender.P2PKHLockingScript().String(), s)
		require.NoErrorf(spec.t, err, "creating utxo for input: %d", i)

		utxo.UnlockingScriptTemplate = spec.sender.P2PKHUnlockingScriptTemplate()

		spec.utxos = append(spec.utxos, utxo)
	}
	spec.sourceTransactions[sourceTX.TxID().String()] = sourceTX

	return spec
}

// WithOPReturn adds an OP_RETURN output to the transaction with the specified data
func (spec *txSpec) WithOPReturn(dataStr string) TransactionSpec {
	data := []byte(dataStr)
	o, err := trx.CreateOpReturnOutput([][]byte{data})
	require.NoError(spec.t, err, "creating op return output")

	spec.outputs = append(spec.outputs, o)

	return spec
}

// WithP2PKHOutput adds a P2PKH output to the transaction with the specified satoshis owned by the recipient
func (spec *txSpec) WithP2PKHOutput(satoshis uint64) TransactionSpec {
	spec.outputs = append(spec.outputs, &trx.TransactionOutput{
		Satoshis:      satoshis,
		LockingScript: spec.recipient.P2PKHLockingScript(),
	})
	return spec
}

// WithOutputScript adds an output to the transaction with the specified satoshis and script
func (spec *txSpec) WithOutputScript(satoshis uint64, script *script.Script) TransactionSpec {
	spec.outputs = append(spec.outputs, &trx.TransactionOutput{
		Satoshis:      satoshis,
		LockingScript: script,
	})
	return spec
}

// ScriptPart is an interface for building script parts
type ScriptPart interface {
	Append(s *script.Script) error
}

// OpCode is an alias for byte to represent an opcode and implements ScriptPart for script building
type OpCode byte

// Append appends the opcode to the script
func (op OpCode) Append(s *script.Script) error {
	return s.AppendOpcodes(script.OpRETURN)
}

// PushData is an alias for []byte to represent push data and implements ScriptPart for script building
type PushData []byte

// Append appends the push data to the script
func (data PushData) Append(s *script.Script) error {
	return s.AppendPushData(data)
}

// WithOutputScriptParts adds an output to the transaction with the specified script parts
func (spec *txSpec) WithOutputScriptParts(parts ...ScriptPart) TransactionSpec {
	s := &script.Script{}
	for _, part := range parts {
		err := part.Append(s)
		require.NoError(spec.t, err, "appending script part")
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
	utxo, err := trx.NewUTXO(spec.fixtureFactory.getNextGrandparentTXID(), 0, spec.sender.P2PKHLockingScript().String(), withFee)
	require.NoError(spec.t, err, "creating utxo for parent tx")

	utxo.UnlockingScriptTemplate = spec.sender.P2PKHUnlockingScriptTemplate()

	err = tx.AddInputsFromUTXOs(utxo)
	require.NoError(spec.t, err, "adding input to parent tx")

	for _, s := range satoshis {
		tx.AddOutput(&trx.TransactionOutput{
			Satoshis:      s,
			LockingScript: spec.sender.P2PKHLockingScript(),
		})
	}
	err = tx.Sign()
	require.NoError(spec.t, err, "signing parent tx")

	// each merkle proof should have a different block height to not collide with each other
	err = tx.AddMerkleProof(trx.NewMerklePath(spec.fixtureFactory.getNextBlockHeight(), [][]*trx.PathElement{{
		&trx.PathElement{
			Hash:   tx.TxID(),
			Offset: 0,
		},
	}}))
	require.NoError(spec.t, err, "adding merkle proof to parent tx")

	return tx
}
