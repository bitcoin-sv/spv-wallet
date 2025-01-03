package record

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
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

	utxosToSpend, trackedOutputs, err := flow.getFromInputs()
	if err != nil {
		return err
	}

	for _, utxo := range utxosToSpend {
		flow.prepareOperationForUserIfNotExist(utxo.UserID)
		flow.subtractSatoshiFromOperation(utxo, utxo.Satoshis)
	}

	flow.spendInputs(trackedOutputs)

	newOutputs, newDataRecords, err := s.processAnnotatedOutputs(tx, &outline.Annotations)
	if err != nil {
		return err
	}

	// TODO: getOutputsForTrackedAddresses

	for _, output := range newOutputs {
		utxo := output.ToUserUTXO()
		if utxo != nil {
			flow.prepareOperationForUserIfNotExist(utxo.UserID)
			flow.addSatoshiToOperation(utxo, utxo.Satoshis)
		}
		flow.createOutputs(output)
	}

	if len(newDataRecords) > 0 {
		flow.prepareOperationForUserIfNotExist(userID)
		flow.txRow.AddData(newDataRecords...)
	}

	if err = flow.verify(); err != nil {
		return err
	}

	if err = flow.broadcast(); err != nil {
		return err
	}

	return flow.save()
}

func (s *Service) processAnnotatedOutputs(tx *trx.Transaction, annotations *transaction.Annotations) ([]database.Output, []*database.Data, error) {
	txID := tx.TxID().String()

	var outputRecords []database.Output
	var dataRecords []*database.Data

	for vout, annotation := range annotations.Outputs {
		if vout >= len(tx.Outputs) {
			return nil, nil, txerrors.ErrAnnotationIndexOutOfRange
		}
		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
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
			outputRecords = append(outputRecords, database.NewDataOutput(txID, voutU32))
		case bucket.BSV:
			//TODO
			s.logger.Warn().Msgf("support for BSV bucket is not implemented yet")
		default:
			s.logger.Warn().Msgf("Unknown annotation bucket %s", annotation.Bucket)
			continue
		}
	}

	return outputRecords, dataRecords, nil
}
