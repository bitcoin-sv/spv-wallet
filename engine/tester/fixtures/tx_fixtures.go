package fixtures

import (
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
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
// parentTX txs contains merkle proof
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
	WithT(t *testing.T) GivenTXSpec
	WithSign() GivenTXSpec
	WithInput(satoshis uint64) GivenTXSpec
	WithOPReturn(dataStr string) GivenTXSpec
	WithOutputFunc(maker func() *trx.TransactionOutput) GivenTXSpec

	TX() *trx.Transaction
	InputUTXO(inputID int) bsv.Outpoint
	ID() string
	BeefHex() string
	RawTXHex() string
	EFHex() string
}

type txSpec struct {
	utxos   []*trx.UTXO
	outputs []*trx.TransactionOutput
	t       *testing.T
	sign    bool

	grandparentTXIndex int
	sourceTransactions map[string]*trx.Transaction
}

// GivenTX creates a new GivenTXSpec for building a MOCK! transaction
func GivenTX() GivenTXSpec {
	return &txSpec{
		grandparentTXIndex: 0,
		sourceTransactions: make(map[string]*trx.Transaction),
	}
}

// WithT sets the testing.T instance to be used for error reporting
func (spec *txSpec) WithT(t *testing.T) GivenTXSpec {
	spec.t = t
	return spec
}

// WithSign enables signing of the transaction (default is false)
func (spec *txSpec) WithSign() GivenTXSpec {
	spec.sign = true
	return spec
}

// WithInput adds an input to the transaction with the specified satoshis
// it automatically creates a parent tx (sourceTX) with P2PKH UTXO with provided satoshis
func (spec *txSpec) WithInput(satoshis uint64) GivenTXSpec {
	sourceTX := spec.makeParentTX(satoshis)
	utxo, err := trx.NewUTXO(sourceTX.TxID().String(), 0, P2PKHLockingScript.String(), satoshis)
	if err != nil {
		spec.error("creating utxo for input", err)
	}
	utxo.UnlockingScriptTemplate = P2PKHUnlockingScriptTemplate

	spec.sourceTransactions[sourceTX.TxID().String()] = sourceTX
	spec.utxos = append(spec.utxos, utxo)

	return spec
}

// WithOPReturn adds an OP_RETURN output to the transaction with the specified data
func (spec *txSpec) WithOPReturn(dataStr string) GivenTXSpec {
	data := []byte(dataStr)
	o, err := trx.CreateOpReturnOutput([][]byte{data})
	if err != nil {
		spec.error("creating op return output", err)
	}

	spec.outputs = append(spec.outputs, o)

	return spec
}

// WithOutputFunc can be used for adding custom outputs to the transaction
func (spec *txSpec) WithOutputFunc(maker func() *trx.TransactionOutput) GivenTXSpec {
	spec.outputs = append(spec.outputs, maker())
	return spec
}

// TX creates a go-sdk transaction from the given spec
func (spec *txSpec) TX() *trx.Transaction {
	tx := trx.NewTransaction()
	if err := tx.AddInputsFromUTXOs(spec.utxos...); err != nil {
		spec.error("adding inputs to tx", err)
	}

	for _, output := range spec.outputs {
		tx.AddOutput(output)
	}

	//TODO consider to source inputs under some condition (spec param)
	for _, input := range tx.Inputs {
		if sourceTX := spec.sourceTransactions[input.SourceTXID.String()]; sourceTX != nil {
			input.SourceTransaction = sourceTX
		}
	}

	if spec.sign {
		if err := tx.Sign(); err != nil {
			spec.error("signing tx", tx, err)
		}
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

// ID returns the transaction ID
func (spec *txSpec) ID() string {
	return spec.TX().TxID().String()
}

// BeefHex returns the BEEF hex of the transaction
func (spec *txSpec) BeefHex() string {
	tx := spec.TX()
	beef, err := tx.BEEFHex()
	if err != nil {
		spec.error("getting beef hex", err)
	}
	return beef
}

// RawTXHex returns the raw hex of the transaction
func (spec *txSpec) RawTXHex() string {
	tx := spec.TX()
	return tx.Hex()
}

// EFHex returns the EF hex of the transaction
func (spec *txSpec) EFHex() string {
	tx := spec.TX()
	ef, err := tx.EFHex()
	if err != nil {
		spec.error("getting ef hex", err)
	}
	return ef
}

func (spec *txSpec) error(args ...any) {
	if spec.t != nil {
		spec.t.Error(args...)
	} else {
		panic(args)
	}
}

func (spec *txSpec) makeParentTX(satoshis uint64) *trx.Transaction {
	tx := trx.NewTransaction()

	withFee := satoshis + 1
	utxo, err := trx.NewUTXO(spec.getNextGrandparentTXID(), 0, P2PKHLockingScript.String(), withFee)
	if err != nil {
		spec.error("creating utxo for source tx", err)
	}
	utxo.UnlockingScriptTemplate = P2PKHUnlockingScriptTemplate

	err = tx.AddInputsFromUTXOs(utxo)
	if err != nil {
		spec.error("adding input to source tx", err)
	}

	tx.AddOutput(&trx.TransactionOutput{
		Satoshis:      satoshis,
		LockingScript: P2PKHLockingScript,
	})
	err = tx.Sign()
	if err != nil {
		spec.error("signing source tx", err)
	}

	err = tx.AddMerkleProof(trx.NewMerklePath(1000, [][]*trx.PathElement{{
		&trx.PathElement{
			Hash:   tx.TxID(),
			Offset: 0,
		},
	}}))
	if err != nil {
		spec.error("adding merkle proof to source tx", err)
	}

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
