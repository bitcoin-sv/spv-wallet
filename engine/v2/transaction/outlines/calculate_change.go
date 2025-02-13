package outlines

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

func calculateChange(inputs annotatedInputs, outputs annotatedOutputs, fee bsv.Satoshis) (bsv.Satoshis, error) {
	satsIn := inputs.totalSatoshis()
	satsOut := outputs.totalSatoshis()

	if satsIn < satsOut+fee {
		// this should never happen
		// UTXO selector should deduce if change output is required and then select enough funds to cover the fee for additional size from the change output
		// NOTE: If user doesn't have enough funds to cover transaction the txerrors.ErrTxOutlineInsufficientFunds is returned on another level
		return 0, spverrors.Wrapf(txerrors.ErrUTXOSelectorInsufficientInputs, "satsIn (%d) are less than satsOut (%d) plus fee (%d)", satsIn, satsOut, fee)
	}

	// keep in mind that those values are uint64 so be careful with underflow
	// that's why we're checking if satsIn is less than satsOut + fee before subtracting
	// don't move this line before the check
	change := satsIn - satsOut - fee

	return change, nil
}
