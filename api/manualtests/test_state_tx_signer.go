package manualtests

import (
	"strconv"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/joomcode/errorx"
	"github.com/samber/lo"
)

var (
	anyonePriv, _ = primitives.PrivateKeyFromBytes([]byte{1})
	anyonePub     = anyonePriv.PubKey()
)

var SigningError = StateError.NewSubtype("signing")

type TxSigner struct {
	state                    *State
	inputs                   map[string]client.ModelsInputAnnotation
	unlockingScriptTemplates map[int]sdk.UnlockingScriptTemplate
}

func NewTxSigner(s *State) *TxSigner {
	return &TxSigner{
		state: s,
	}
}

func (t *TxSigner) UnlockToHex(format string, hex string) (string, error) {
	tx, err := t.Unlock(format, hex)
	if err != nil {
		return "", err
	}
	switch format {
	case "RAW", "raw":
		return tx.Hex(), nil
	case "BEEF", "beef":
		beef, err := tx.BEEFHex()
		if err != nil {
			return "", SigningError.Wrap(err, "failed to get BEEF hex from signed transaction")
		}
		return beef, nil
	default:
		return "", SigningError.New("unknown format: %s", format)
	}
}

func (t *TxSigner) Unlock(format string, hex string) (*sdk.Transaction, error) {
	err := t.prepareUnlockingScriptsFromAnnotations()
	if err != nil {
		return nil, err
	}

	switch format {
	case "RAW", "raw":
		return t.UnlockRawTx(hex)
	case "BEEF", "beef":
		return t.UnlockBeefTx(hex)
	default:
		return nil, SigningError.New("unknown format: %s", format)
	}
}

func (t *TxSigner) UsingAnnotations(inputs map[string]client.ModelsInputAnnotation) *TxSigner {
	t.inputs = inputs
	return t
}

func (t *TxSigner) prepareUnlockingScriptsFromAnnotations() (err error) {
	t.unlockingScriptTemplates = make(map[int]sdk.UnlockingScriptTemplate, len(t.inputs))
	for k, input := range t.inputs {
		idx, err := strconv.Atoi(k)
		if err != nil {
			return SigningError.Wrap(err, "failed to convert input key %s to int", k)
		}

		t.unlockingScriptTemplates[idx], err = t.prepareUnlockingScriptTemplate(input)
		if err != nil {
			t.unlockingScriptTemplates = nil
			return errorx.Decorate(err, "failed to prepare unlocking script template for input %d", idx)
		}
	}
	return nil
}

func (t *TxSigner) prepareUnlockingScriptTemplate(input client.ModelsInputAnnotation) (sdk.UnlockingScriptTemplate, error) {
	instructions, err := input.CustomInstructions.AsModelsSPVWalletCustomInstructions()
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to get custom instructions from input annotations, maybe it's not spv-wallet custom instruction?")
	}

	key := t.state.CurrentUser().GetPrivateKey()
	for i, ci := range instructions {
		switch ci.Type {
		case "type42":
			key, err = key.DeriveChild(anyonePub, ci.Instruction)
			if err != nil {
				return nil, SigningError.Wrap(err, "failed to derive key for type42 instruction (index = %d) with invoiceNumber %s", i, ci.Instruction)
			}
		default:
			return nil, SigningError.New("unknown (or not implemented) instruction (index = %d) type: %s with instruction %s", i, ci.Type, ci.Instruction)
		}
	}

	unlockingScript, err := p2pkh.Unlock(key, lo.ToPtr(sighash.AllForkID))
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to create unlocking script for input %s and key %s", input, key.Wif())
	}
	return unlockingScript, nil
}

func (t *TxSigner) UnlockBeefTx(hex string) (*sdk.Transaction, error) {
	tx, err := sdk.NewTransactionFromBEEFHex(hex)
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to create transaction from BEEF hex %s", hex)
	}
	err = t.UnlockTx(tx)
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to sign BEEF transaction")
	}
	return tx, nil
}

func (t *TxSigner) UnlockRawTx(hex string) (*sdk.Transaction, error) {
	tx, err := sdk.NewTransactionFromHex(hex)
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to create transaction from hex %s", hex)
	}
	err = t.UnlockTx(tx)
	if err != nil {
		return nil, SigningError.Wrap(err, "failed to sign RAW transaction")
	}
	return tx, nil
}

func (t *TxSigner) UnlockTx(tx *sdk.Transaction) error {
	if t.inputs == nil {
		return SigningError.New("inputs annotation not set")
	}

	for idx, unlockingScriptTemplate := range t.unlockingScriptTemplates {
		input := tx.InputIdx(idx)
		if input == nil {
			return SigningError.New("input %d not found in transaction but has unlocking instructions", idx)
		}
		input.UnlockingScriptTemplate = unlockingScriptTemplate
	}

	err := tx.Sign()
	if err != nil {
		return SigningError.Wrap(err, "failed to sign transaction")
	}
	return nil
}
