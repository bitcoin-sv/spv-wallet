package record

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// RecordTransactionOutline will validate, broadcast and save a transaction outline
func (s *Service) RecordTransactionOutline(ctx context.Context, outline *outlines.Transaction) error {
	tx, err := trx.NewTransactionFromBEEFHex(outline.BEEF)
	if err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	}

	if ok, err := spv.VerifyScripts(tx); err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	} else if !ok {
		return txerrors.ErrTxValidation
	}

	utxos, err := s.getTrackedUTXOsFromInputs(ctx, tx)
	if err != nil {
		return err
	}

	txID := tx.TxID().String()
	for _, utxo := range utxos {
		utxo.Spend(txID)
	}

	var newOutputs []*database.Output
	var newDataRecords []*database.Data
	if outline.Annotations != nil {
		newOutputs, newDataRecords, err = s.processAnnotatedOutputs(tx, *outline.Annotations)
		if err != nil {
			return err
		}
	}

	if err = s.broadcaster.Broadcast(ctx, tx); err != nil {
		return txerrors.ErrTxBroadcast.Wrap(err)
	}

	txRecord := database.Transaction{
		ID:       txID,
		TxStatus: database.TxStatusBroadcasted,
	}

	upsertOutputs := append(newOutputs, utxos...) //newly created outputs and spent utxos
	err = s.repo.SaveTX(ctx, &txRecord, upsertOutputs, newDataRecords)
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}

	return nil
}

func (s *Service) getTrackedUTXOsFromInputs(ctx context.Context, tx *trx.Transaction) ([]*database.Output, error) {
	outpoints := func(yield func(outpoint bsv.Outpoint) bool) {
		for _, input := range tx.Inputs {
			yield(bsv.Outpoint{
				TxID: input.SourceTXID.String(),
				Vout: input.SourceTxOutIndex,
			})
		}
	}
	storedUTXOs, err := s.repo.GetOutputs(ctx, outpoints)
	if err != nil {
		return nil, txerrors.ErrGettingOutputs.Wrap(err)
	}

	for _, utxo := range storedUTXOs {
		if utxo.IsSpent() {
			return nil, txerrors.ErrUTXOSpent
		}
	}

	return storedUTXOs, nil
}

func (s *Service) processAnnotatedOutputs(tx *trx.Transaction, annotations transaction.Annotations) ([]*database.Output, []*database.Data, error) {
	txID := tx.TxID().String()

	var outputRecords []*database.Output
	var dataRecords []*database.Data

	for vout, annotation := range annotations.Outputs {
		if vout >= len(tx.Outputs) {
			return nil, nil, txerrors.ErrAnnotationIndexOutOfRange
		}
		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, nil, spverrors.Wrapf(err, "Vout value exceeds max uint32 range")
		}
		lockingScript := tx.Outputs[vout].LockingScript

		switch annotation.Bucket {
		case bucket.Data:
			data, err := getDataFromOpReturn(lockingScript)
			if err != nil {
				return nil, nil, err
			}
			dataRecords = append(dataRecords, &database.Data{
				TxID: txID,
				Vout: voutU32,
				Blob: data,
			})
		case bucket.BSV: //TODO
		default:
			s.logger.Warn().Msgf("Unknown annotation bucket %s", annotation.Bucket)
			continue
		}
	}

	return outputRecords, dataRecords, nil
}
