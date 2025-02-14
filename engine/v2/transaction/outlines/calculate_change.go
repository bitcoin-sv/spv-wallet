package outlines

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func calculateChange(ctx *evaluationContext, inputs annotatedInputs, outputs annotatedOutputs, feeUnit bsv.FeeUnit) (annotatedOutputs, error) {
	satsIn := inputs.totalSatoshis()
	satsOut := outputs.totalSatoshis()

	fee := calculateFee(inputs, outputs, feeUnit)
	if satsIn < satsOut+fee {
		// this should never happen
		// UTXO selector should provide enough funds to cover outputs and fee
		// NOTE: If user doesn't have enough funds to cover transaction the txerrors.ErrTxOutlineInsufficientFunds is returned on another level
		return nil, spverrors.Wrapf(txerrors.ErrUTXOSelectorInsufficientInputs, "satsIn (%d) are less than satsOut (%d) plus fee (%d)", satsIn, satsOut, fee)
	}

	// keep in mind that those values are uint64 so be careful with underflow
	// that's why we're checking if satsIn is less than satsOut + fee before subtracting
	// don't move this line before the check
	change := satsIn - satsOut - fee

	if change == 0 {
		return outputs, nil
	}

	userPubKey, err := ctx.UserPubKey()
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user public key")
	}

	lockingScript, customInstructions, err := lockingScriptForChangeOutput(userPubKey)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create locking script for change output")
	}
	changeOutput := &annotatedOutput{
		OutputAnnotation: &transaction.OutputAnnotation{
			Bucket:             bucket.BSV,
			CustomInstructions: &customInstructions,
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
			// Addition of change output increased fee, and decreases change to 0
			// returning outputs with higher fee and no change
			return outputs, nil
		}

		changeOutput.TransactionOutput.Satoshis -= feeDiff
	}

	return outputsWithChange, nil
}

func lockingScriptForChangeOutput(pubKey *primitives.PublicKey) (*script.Script, bsv.CustomInstructions, error) {
	dest, err := type42.NewDestinationWithRandomReference(pubKey)
	if err != nil {
		return nil, nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	address, err := script.NewAddressFromPublicKey(dest.PubKey, true)
	if err != nil {
		return nil, nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	lockingScript, err := p2pkh.Lock(address)
	if err != nil {
		return nil, nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	customInstructions := bsv.CustomInstructions{
		{
			Type:        "type42",
			Instruction: dest.DerivationKey,
		},
	}

	return lockingScript, customInstructions, nil
}
