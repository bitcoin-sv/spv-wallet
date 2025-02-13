package outlines

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func calculateChange(ctx *evaluationContext, inputs annotatedInputs, outputs annotatedOutputs, feeUnit bsv.FeeUnit) (annotatedOutputs, bsv.Satoshis, error) {
	satsIn := inputs.totalSatoshis()
	satsOut := outputs.totalSatoshis()

	fee := calculateFee(inputs, outputs, feeUnit)
	if satsIn < satsOut+fee {
		// this should never happen
		// UTXO selector should deduce if change output is required and then select enough funds to cover the fee for additional size from the change output
		// NOTE: If user doesn't have enough funds to cover transaction the txerrors.ErrTxOutlineInsufficientFunds is returned on another level
		return nil, 0, spverrors.Wrapf(txerrors.ErrUTXOSelectorInsufficientInputs, "satsIn (%d) are less than satsOut (%d) plus fee (%d)", satsIn, satsOut, fee)
	}

	// keep in mind that those values are uint64 so be careful with underflow
	// that's why we're checking if satsIn is less than satsOut + fee before subtracting
	// don't move this line before the check
	change := satsIn - satsOut - fee

	if change == 0 {
		return outputs, 0, nil
	}

	userPubKey, err := ctx.UserPubKey()
	if err != nil {
		return nil, 0, spverrors.Wrapf(err, "failed to get user public key")
	}

	lockingScript, err := lockingScriptForChangeOutput(userPubKey)
	if err != nil {
		return nil, 0, spverrors.Wrapf(err, "failed to create locking script for change output")
	}
	changeOutput := &annotatedOutput{
		OutputAnnotation: &transaction.OutputAnnotation{
			Bucket:             bucket.BSV,
			CustomInstructions: nil, //fixme: add custom instructions
		},
		TransactionOutput: &sdk.TransactionOutput{
			LockingScript: lockingScript,
			Satoshis:      uint64(change),
		},
	}
	outputsWithChange := append(outputs, changeOutput)

	feeForTxWithChange := calculateFee(inputs, outputsWithChange, feeUnit)
	if feeForTxWithChange > fee {
		feeDiff := uint64(feeForTxWithChange - fee)
		if feeDiff >= changeOutput.TransactionOutput.Satoshis {
			return outputs, 0, nil
		}
		changeOutput.TransactionOutput.Satoshis -= feeDiff
	}

	return outputsWithChange, change, nil
}

func lockingScriptForChangeOutput(pubKey *primitives.PublicKey) (*script.Script, error) {
	addr, err := script.NewAddressFromPublicKey(pubKey, true)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create change output address")
	}

	lockingScript, err := p2pkh.Lock(addr)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create locking script for change output")
	}

	return lockingScript, nil
}
