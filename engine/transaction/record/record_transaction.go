package record

import (
	"context"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// RecordPaymailTransaction will validate, broadcast and save paymail transaction
func (s *Service) RecordPaymailTransaction(ctx context.Context, tx *trx.Transaction, senderPaymail, receiverPaymail string) error {
	flow := newTxFlow(ctx, s, tx)

	trackedOutputs, err := flow.getFromInputs()
	if err != nil {
		return err
	}

	for _, utxo := range trackedOutputs {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", receiverPaymail)
		if len(flow.operations) > 1 {
			return spverrors.Newf("paymail transaction with multiple senders is not supported")
		}
		operation.subtract(utxo.Satoshis)
	}

	flow.spendInputs(trackedOutputs)

	p2pkhOutputs, err := flow.findRelevantP2PKHOutputs()
	if err != nil {
		return err
	}

	for outputData := range p2pkhOutputs {
		operation := flow.operationOfUser(outputData.userID, "incoming", senderPaymail)
		if len(flow.operations) > 2 {
			return spverrors.Newf("paymail transaction with multiple receivers is not supported")
		}
		operation.add(outputData.satoshis)
		flow.createP2PKHOutput(outputData)
	}

	if err = flow.verify(); err != nil {
		return err
	}

	if err = flow.broadcast(); err != nil {
		return err
	}

	return flow.save()
}
