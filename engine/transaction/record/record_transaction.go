package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// RecordPaymailTransaction will validate, broadcast and save paymail transaction
func (s *Service) RecordPaymailTransaction(ctx context.Context, tx *trx.Transaction, senderPaymail, receiverPaymail string) error {
	flow := newTxFlow(ctx, s, tx)

	utxosToSpend, trackedOutputs, err := flow.getFromInputs()
	if err != nil {
		return err
	}

	for _, utxo := range utxosToSpend {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", receiverPaymail)
		if len(flow.operations) > 1 {
			return spverrors.Newf("paymail transaction with multiple senders is not supported")
		}
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
			operation := flow.operationOfUser(utxo.UserID, "incoming", senderPaymail)
			if len(flow.operations) > 2 {
				return spverrors.Newf("paymail transaction with multiple receivers is not supported")
			}
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
