package utils

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
)

// GetUnlockingScript will generate an unlocking script
func GetUnlockingScript(tx *bt.Tx, inputIndex uint32, privateKey *bec.PrivateKey) (*bscript.Script, error) {
	sigHashFlags := sighash.AllForkID

	sigHash, err := tx.CalcInputSignatureHash(inputIndex, sigHashFlags)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to calculate signature hash")
	}

	var sig *bec.Signature
	if sig, err = privateKey.Sign(sigHash); err != nil {
		return nil, spverrors.Wrapf(err, "failed to sign transaction")
	}

	pubKey := privateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	var script *bscript.Script
	if script, err = bscript.NewP2PKHUnlockingScript(pubKey, signature, sigHashFlags); err != nil {
		return nil, spverrors.Wrapf(err, "failed to create unlocking script")
	}

	return script, nil
}
