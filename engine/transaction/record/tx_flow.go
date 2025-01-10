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

	operations map[string]*operationWrapper
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

		operations: map[string]*operationWrapper{},
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

func (f *txFlow) operationOfUser(userID string) *operationWrapper {
	if _, ok := f.operations[userID]; !ok {
		f.operations[userID] = &operationWrapper{
			entity: &database.Operation{
				UserID: userID,

				Transaction: f.txRow,
				Value:       0,
			},
		}
	}
	return f.operations[userID]
}

func (f *txFlow) spendInputs(trackedOutputs []*database.TrackedOutput) {
	f.txRow.AddInputs(trackedOutputs...)
}

func (f *txFlow) createOutputs(outputs ...database.Output) {
	f.txRow.AddOutputs(outputs...)
}

func (f *txFlow) getOutputsForTrackedAddresses() (iter.Seq[database.Output], error) {
	relevantOutputs := map[string]uint32{} // address -> vout
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

		voutU32, err := conv.IntToUint32(vout)
		if err != nil {
			f.service.logger.Warn().Err(err).Msg("failed to convert vout to uint32")
			continue
		}

		relevantOutputs[address.AddressString] = voutU32
	}

	rows, err := f.service.repo.GetAddresses(f.ctx, maps.Keys(relevantOutputs))
	if err != nil {
		return nil, txerrors.ErrGettingAddresses.Wrap(err)
	}

	return func(yield func(database.Output) bool) {
		for _, row := range rows {
			vout, ok := relevantOutputs[row.Address]
			if !ok {
				f.service.logger.Warn().Str("address", row.Address).Msg("Got not relevant address from database")
				continue
			}
			yield(database.NewP2PKHOutput(
				f.txID,
				vout,
				row.UserID,
				bsv.Satoshis(f.tx.Outputs[vout].Satoshis)),
			)
		}
	}, nil
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
	err := f.service.repo.SaveOperations(f.ctx, toOperationEntities(maps.Values(f.operations)))
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}
	return nil
}
