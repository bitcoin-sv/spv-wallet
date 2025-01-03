package record

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
)

func (s *Service) RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error {
	if verifyScripts {
		// TODO: Check if in case of not-veryfying-scripts we accidentally allow removing UserUTXOs from the database
		// NOTE: When we want to record "RawTX" we cannot verify scripts
		if ok, err := spv.VerifyScripts(tx); err != nil {
			return txerrors.ErrTxValidation.Wrap(err)
		} else if !ok {
			return txerrors.ErrTxValidation
		}
	}

	utxos, err := s.getTrackedUTXOsFromInputs(ctx, tx)
	if err != nil {
		return err
	}

	newOutputs, err := s.getOutputsForTrackedAddresses(ctx, tx)
	if err != nil {
		return err
	}

	if _, err = s.broadcaster.Broadcast(ctx, tx); err != nil {
		return txerrors.ErrTxBroadcast.Wrap(err)
	}
	// TODO: handle TXInfo returned from Broadcast (SPV-1157)

	txID := tx.TxID().String()

	txRow := database.TrackedTransaction{
		ID:       txID,
		TxStatus: database.TxStatusBroadcasted,
	}
	txRow.AddInputs(utxos...)
	txRow.AddOutputs(newOutputs...)

	err = s.repo.SaveTX(ctx, &txRow)
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}

	return nil
}

func (s *Service) getOutputsForTrackedAddresses(ctx context.Context, tx *trx.Transaction) ([]database.Output, error) {
	var trackedOutputs []database.Output
	for vout, output := range tx.Outputs {
		lockingScript := output.LockingScript
		if !lockingScript.IsP2PKH() {
			continue
		}
		address, err := lockingScript.Address()
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to get address from locking script")
			continue
		}

		addressRow, err := s.repo.CheckAddress(ctx, address.AddressString)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to check address")
			continue
		}
		if addressRow == nil || addressRow.User == nil {
			s.logger.Debug().Str("address", address.AddressString).Msg("address is not tracked")
			continue
		}

		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
		}

		trackedOutputs = append(trackedOutputs, database.NewP2PKHOutput(
			tx.TxID().String(),
			voutU32,
			addressRow.User.ID,
			bsv.Satoshis(output.Satoshis)),
		)
	}
	return trackedOutputs, nil
}
