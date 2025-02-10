package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
)

// RecordPaymailTransaction will validate, broadcast and save paymail transaction
func (s *Service) RecordPaymailTransaction(ctx context.Context, tx *trx.Transaction, senderPaymail, receiverPaymail string) error {
	flow := newTxFlow(ctx, s, tx)

	trackedOutputs, err := flow.processInputs()
	if err != nil {
		return err
	}

	for _, utxo := range trackedOutputs {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", receiverPaymail)
		if len(flow.operations) > 1 {
			return spverrors.Newf("paymail transaction with multiple senders is not supported")
		}
		operation.Subtract(utxo.Satoshis)
	}

	p2pkhOutputs, err := flow.findRelevantP2PKHOutputs()
	if err != nil {
		return err
	}

	for outputData := range p2pkhOutputs {
		operation := flow.operationOfUser(outputData.UserID, "incoming", senderPaymail)
		if len(flow.operations) > 2 {
			return txerrors.ErrMultiPaymailRecipientsNotSupported
		}
		operation.Add(outputData.Satoshis)
		flow.addOutputs(outputData)
	}

	if err = flow.verify(); err != nil {
		return err
	}

	if err = flow.broadcast(); err != nil {
		return err
	}

	return flow.save()
}
