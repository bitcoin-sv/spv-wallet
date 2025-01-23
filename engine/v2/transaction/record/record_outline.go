package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// RecordTransactionOutline will validate, broadcast and save a transaction outline
func (s *Service) RecordTransactionOutline(ctx context.Context, userID string, outline *outlines.Transaction) (*txmodels.RecordedOutline, error) {
	tx, err := trx.NewTransactionFromBEEFHex(outline.BEEF)
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

	newDataRecords, err := processDataOutputs(tx, userID, &outline.Annotations)
	if err != nil {
		return nil, err
	}

	// TODO: getOutputsForTrackedAddresses
	// TODO: process Paymail Annotations

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
