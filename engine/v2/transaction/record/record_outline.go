package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// RecordTransactionOutline will validate, broadcast and save a transaction outline
func (s *Service) RecordTransactionOutline(ctx context.Context, userID string, outline *outlines.Transaction) error {
	tx, err := trx.NewTransactionFromBEEFHex(outline.BEEF)
	if err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	}

	flow := newTxFlow(ctx, s, tx)
	if err = flow.verifyScripts(); err != nil {
		return err
	}

	trackedOutputs, err := flow.getFromInputs()
	if err != nil {
		return err
	}

	for _, utxo := range trackedOutputs {
		operation := flow.operationOfUser(utxo.UserID, "outgoing", "")
		operation.subtract(utxo.Satoshis)
	}

	flow.spendInputs(trackedOutputs)

	newDataRecords, err := s.processDataOutputs(tx, &outline.Annotations)
	if err != nil {
		return err
	}

	// TODO: getOutputsForTrackedAddresses
	// TODO: process Paymail Annotations

	if len(newDataRecords) > 0 {
		_ = flow.operationOfUser(userID, "data", "")
		flow.createDataOutputs(userID, newDataRecords...)
	}

	if err = flow.verify(); err != nil {
		return err
	}

	if err = flow.broadcast(); err != nil {
		return err
	}

	return flow.save()
}

func (s *Service) processDataOutputs(tx *trx.Transaction, annotations *transaction.Annotations) ([]*database.Data, error) {
	txID := tx.TxID().String()

	var dataRecords []*database.Data //nolint: prealloc

	for vout, annotation := range annotations.Outputs {
		if vout >= len(tx.Outputs) {
			return nil, txerrors.ErrAnnotationIndexOutOfRange
		}
		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
		}
		lockingScript := tx.Outputs[vout].LockingScript

		if annotation.Bucket != bucket.Data {
			continue
		}

		data, err := getDataFromOpReturn(lockingScript)
		if err != nil {
			return nil, err
		}
		dataRecords = append(dataRecords, &database.Data{
			TxID: txID,
			Vout: voutU32,
			Blob: data,
		})
	}

	return dataRecords, nil
}
