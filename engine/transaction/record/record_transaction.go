package record

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/conv"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
)

func (s *Service) RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error {
	if verifyScripts {
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

func (s *Service) getOutputsForTrackedAddresses(ctx context.Context, tx *trx.Transaction) ([]*database.TrackedOutput, error) {
	var trackedOutputs []*database.TrackedOutput
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

		tracked, err := s.repo.CheckAddress(ctx, address.AddressString)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to check address")
			continue
		}
		if !tracked {
			s.logger.Debug().Str("address", address.AddressString).Msg("address not tracked")
			continue
		}

		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
		}

		trackedOutputs = append(trackedOutputs, &database.TrackedOutput{
			TxID: tx.TxID().String(),
			Vout: voutU32,
		})
	}
	return trackedOutputs, nil
}