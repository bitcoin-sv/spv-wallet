package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

// RecordTransaction will validate, broadcast and save a transaction
func (s *Service) RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error {
	flow := newTxFlow(ctx, s, tx)

	if verifyScripts {
		// TODO: Check if in case of not-veryfying-scripts we accidentally allow removing UserUTXOs from the database
		// NOTE: When we want to record "RawTX" we cannot verify scripts
		if err := flow.verifyScripts(); err != nil {
			return err
		}
	}

	utxosToSpend, trackedOutputs, err := flow.getFromInputs()
	if err != nil {
		return err
	}

	for _, utxo := range utxosToSpend {
		operation := flow.operationOfUser(utxo.UserID)
		operation.subtract(utxo.Satoshis)
	}

	flow.spendInputs(trackedOutputs)

	newOutputs, err := flow.getOutputsForTrackedAddresses()
	if err != nil {
		return err
	}

	for output := range newOutputs {
		utxo := output.ToUserUTXO()
		if utxo != nil {
			operation := flow.operationOfUser(utxo.UserID)
			operation.add(utxo.Satoshis)
		}
		flow.createOutputs(output)
	}

	if err = flow.verify(); err != nil {
		return err
	}

	if err = flow.broadcast(); err != nil {
		return err
	}

	return flow.save()
}
