package record

import (
	"context"
	"github.com/bitcoin-sv/go-sdk/script"
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

func (s *Service) RecordTransactionOutline(ctx context.Context, outline *outlines.Transaction) error {
	tx, err := trx.NewTransactionFromBEEFHex(outline.BEEF)
	if err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	}

	txID := tx.TxID().String()

	if ok, err := spv.VerifyScripts(tx); err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	} else if !ok {
		return txerrors.ErrTxValidation
	}

	utxos, err := s.getTrackedUTXOsFromInputs(ctx, tx)
	if err != nil {
		return err
	}
	for _, utxo := range utxos {
		utxo.Spend(txID)
	}

	var outputRecords []database.Output
	var dataRecords []database.Data
	if outline.Annotations != nil {
		outputRecords, dataRecords, err = s.processAnnotatedOutputs(tx, *outline.Annotations)
		if err != nil {
			return err
		}
	}

	txRecord := database.Transaction{
		ID:       txID,
		TxStatus: database.TxStatusCreated,
	}

	if err = s.broadcaster.Broadcast(ctx, tx); err != nil {
		return txerrors.ErrTxBroadcast.Wrap(err)
	} else {
		txRecord.TxStatus = database.TxStatusBroadcasted
	}

	err = s.repo.SaveTX(ctx, &txRecord, outputRecords, dataRecords)
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}

	return nil
}

func (s *Service) getTrackedUTXOsFromInputs(ctx context.Context, tx *trx.Transaction) ([]database.Output, error) {
	txID := tx.TxID().String()
	outpoints := func(yield func(outpoint bsv.Outpoint) bool) {
		for _, input := range tx.Inputs {
			yield(bsv.Outpoint{
				TxID: txID,
				Vout: input.SourceTxOutIndex,
			})
		}
	}
	storedUTXOs, err := s.repo.GetOutputs(ctx, outpoints)
	if err != nil {
		return nil, err //TODO wrap
	}

	for _, utxo := range storedUTXOs {
		if utxo.IsSpent() {
			return nil, txerrors.ErrUTXOSpent
		}
	}

	return storedUTXOs, nil
}

func (s *Service) getDataFromOpReturn(lockingScript *script.Script) ([]byte, error) {
	if !lockingScript.IsData() {
		return nil, spverrors.Newf("Script is not a data output")
	}

	chunks, err := lockingScript.Chunks()
	if err != nil {
		return nil, txerrors.ErrParsingScript.Wrap(err)
	}

	startIndex := 2
	if chunks[0].Op == script.OpRETURN {
		startIndex = 1
	}

	var d [][]byte
	for _, chunk := range chunks[startIndex:] {
		if chunk.Op > script.OpPUSHDATA4 {
			return nil, spverrors.Newf("Could not find OP_RETURN data")
		}
		d = append(d, chunk.Data)
	}

	return d[0], nil
}

func (s *Service) processAnnotatedOutputs(tx *trx.Transaction, annotations transaction.Annotations) ([]database.Output, []database.Data, error) {
	txID := tx.TxID().String()

	var outputRecords []database.Output
	var dataRecords []database.Data

	for vout, annotation := range annotations.Outputs {
		if vout >= len(tx.Outputs) {
			s.logger.Warn().Msgf("Annotation's output index %d is out of range", vout)
			continue
		}
		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, nil, spverrors.Wrapf(err, "Vout value exceeds max uint32 range")
		}
		lockingScript := tx.Outputs[vout].LockingScript

		switch annotation.Bucket {
		case bucket.Data:
			data, err := s.getDataFromOpReturn(lockingScript)
			if err != nil {
				return nil, nil, err
			}
			dataRecords = append(dataRecords, database.Data{
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
