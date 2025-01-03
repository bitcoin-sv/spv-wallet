package record

import (
	"context"
	"iter"
	"maps"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type txFlow struct {
	ctx     context.Context
	service *Service

	tx    *trx.Transaction
	txRow *database.TrackedTransaction
	txID  string

	operations map[string]*database.Operation
}

func newTxFlow(ctx context.Context, service *Service, tx *trx.Transaction) *txFlow {
	txID := tx.TxID().String()
	return &txFlow{
		ctx:     ctx,
		service: service,

		tx:   tx,
		txID: txID,
		txRow: &database.TrackedTransaction{
			ID:       txID,
			TxStatus: database.TxStatusCreated,
		},

		operations: map[string]*database.Operation{},
	}
}

func (f *txFlow) verifyScripts() error {
	if ok, err := spv.VerifyScripts(f.tx); err != nil {
		return txerrors.ErrTxValidation.Wrap(err)
	} else if !ok {
		return txerrors.ErrTxValidation
	}
	return nil
}

func (f *txFlow) getFromInputs() ([]*database.UserUtxos, []*database.TrackedOutput, error) {
	outpoints := func(yield func(outpoint bsv.Outpoint) bool) {
		for _, input := range f.tx.Inputs {
			yield(bsv.Outpoint{
				TxID: input.SourceTXID.String(),
				Vout: input.SourceTxOutIndex,
			})
		}
	}
	utxos, trackedOutputs, err := f.service.repo.GetOutputs(f.ctx, outpoints)
	if err != nil {
		return nil, nil, txerrors.ErrGettingOutputs.Wrap(err)
	}

	for _, output := range trackedOutputs {
		if output.IsSpent() {
			return nil, nil, txerrors.ErrUTXOSpent
		}
	}

	return utxos, trackedOutputs, nil
}

func (f *txFlow) prepareOperationForUserIfNotExist(userID string) {
	if _, ok := f.operations[userID]; !ok {
		f.operations[userID] = &database.Operation{
			UserID: userID,

			Transaction: f.txRow,
			Value:       0,
		}
	}
}

func (f *txFlow) addSatoshiToOperation(utxo *database.UserUtxos, satoshi uint64) {
	signedSatoshi, err := conv.Uint64ToInt64(satoshi)
	if err != nil {
		panic(err)
	}
	f.operations[utxo.UserID].Value = f.operations[utxo.UserID].Value + signedSatoshi
}

func (f *txFlow) subtractSatoshiFromOperation(utxo *database.UserUtxos, satoshi uint64) {
	signedSatoshi, err := conv.Uint64ToInt64(satoshi)
	if err != nil {
		panic(err)
	}
	f.operations[utxo.UserID].Value = f.operations[utxo.UserID].Value - signedSatoshi
}

func (f *txFlow) spendInputs(trackedOutputs []*database.TrackedOutput) {
	f.txRow.AddInputs(trackedOutputs...)
}

func (f *txFlow) createOutputs(outputs ...database.Output) {
	f.txRow.AddOutputs(outputs...)
}

func (f *txFlow) getOutputsForTrackedAddresses() iter.Seq[database.Output] {
	return func(yield func(database.Output) bool) {
		for vout, output := range f.tx.Outputs {
			lockingScript := output.LockingScript
			if !lockingScript.IsP2PKH() {
				continue
			}
			address, err := lockingScript.Address()
			if err != nil {
				f.service.logger.Warn().Err(err).Msg("failed to get address from locking script")
				continue
			}

			addressRow, err := f.service.repo.CheckAddress(f.ctx, address.AddressString)
			if err != nil {
				f.service.logger.Warn().Err(err).Msg("failed to check address")
				continue
			}
			if addressRow == nil || addressRow.UserID == "" {
				f.service.logger.Debug().Str("address", address.AddressString).Msg("address is not tracked")
				continue
			}

			voutU32, err := conv.IntToUint32(vout)
			if err != nil {
				f.service.logger.Warn().Err(err).Msg("failed to convert vout to uint32")
				continue
			}

			yield(database.NewP2PKHOutput(
				f.txID,
				voutU32,
				addressRow.UserID,
				bsv.Satoshis(output.Satoshis)),
			)
		}
	}
}

func (f *txFlow) verify() error {
	if len(f.operations) == 0 {
		return txerrors.ErrNoOperations
	}
	return nil
}

func (f *txFlow) broadcast() error {
	if _, err := f.service.broadcaster.Broadcast(f.ctx, f.tx); err != nil {
		return txerrors.ErrTxBroadcast.Wrap(err)
	}
	return nil
}

func (f *txFlow) save() error {
	err := f.service.repo.SaveOperations(f.ctx, maps.Values(f.operations))
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}
	return nil
}
