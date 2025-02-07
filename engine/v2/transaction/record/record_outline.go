package record

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

// RecordTransactionOutline will validate, broadcast and save a transaction outline
func (s *Service) RecordTransactionOutline(ctx context.Context, userID string, outline *outlines.Transaction) (*txmodels.RecordedOutline, error) {
	if outline.Hex.IsRawTx() {
		return nil, spverrors.Newf("not implemented recording outline with raw transaction")
	}

	tx, err := outline.Hex.ToBEEFTransaction()
	if err != nil {
		return nil, txerrors.ErrTxValidation.Wrap(err)
	}

	flow := newTxFlow(ctx, s, tx)
	if err = flow.verifyScripts(); err != nil {
		return nil, err
	}

	sender := ""
	receiver := ""

	pmInfo, err := processPaymailOutputs(outline.Annotations)
	if err != nil {
		return nil, err
	} else if pmInfo != nil {
		sender = pmInfo.Sender
		receiver = pmInfo.Receiver
	}

	trackedOutputs, err := flow.processInputs()
	if err != nil {
		return nil, err
	}

	for _, utxo := range trackedOutputs {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", receiver)
		operation.Subtract(utxo.Satoshis)
	}

	// getting all outputs that matches user's addresses from the database
	p2pkhOutputs, err := flow.findRelevantP2PKHOutputs()
	if err != nil {
		return nil, err
	}

	for outputData := range p2pkhOutputs {
		if pmInfo != nil && pmInfo.hasVOut(outputData.Vout) {
			// If the output which matches an address obtained from our database,
			// is marked as paymail output in the annotation,
			// it means that we don't have to make paymail-p2p-notification because it is internal.
			pmInfo = nil
		}
		operation := flow.operationOfUser(outputData.UserID, "incoming", sender)
		operation.Add(outputData.Satoshis)
		flow.addOutputs(outputData)
	}

	newDataRecords, err := processDataOutputs(tx, userID, &outline.Annotations)
	if err != nil {
		return nil, err
	}

	if len(newDataRecords) > 0 {
		_ = flow.operationOfUser(userID, "data", sender)
		flow.addOutputs(newDataRecords...)
	}

	if err = flow.verify(); err != nil {
		return nil, err
	}

	if pmInfo != nil {
		if err = flow.notifyPaymailExternalRecipient(*pmInfo); err != nil {
			return nil, err
		}
	}

	if err = flow.broadcast(); err != nil {
		return nil, err
	}

	if err = flow.save(); err != nil {
		return nil, err
	}

	return &txmodels.RecordedOutline{
		TxID: tx.TxID().String(),
	}, nil
}
