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

	trackedOutputs, err := flow.processInputs()
	if err != nil {
		return nil, err
	}

	for _, utxo := range trackedOutputs {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", "")
		operation.Subtract(utxo.Satoshis)
	}

	p2pkhOutputs, err := flow.findRelevantP2PKHOutputs()
	if err != nil {
		return nil, err
	}

	for outputData := range p2pkhOutputs {
		operation := flow.operationOfUser(outputData.UserID, "incoming", "")
		if len(flow.operations) > 2 {
			return nil, spverrors.Newf("paymail transaction with multiple receivers is not supported")
		}
		operation.Add(outputData.Satoshis)
		flow.addOutputs(outputData)
	}

	newDataRecords, err := processDataOutputs(tx, userID, &outline.Annotations)
	if err != nil {
		return nil, err
	}

	if len(newDataRecords) > 0 {
		_ = flow.operationOfUser(userID, "data", "")
		flow.addOutputs(newDataRecords...)
	}

	if err = flow.verify(); err != nil {
		return nil, err
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
