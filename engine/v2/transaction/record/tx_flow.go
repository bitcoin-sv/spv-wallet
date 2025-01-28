package record

import (
	"context"
	"iter"
	"maps"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type txFlow struct {
	ctx     context.Context
	service *Service

	tx    *trx.Transaction
	txRow txmodels.NewTransaction
	txID  string

	operations map[string]*txmodels.NewOperation
}

func newTxFlow(ctx context.Context, service *Service, tx *trx.Transaction) *txFlow {
	txID := tx.TxID().String()
	return &txFlow{
		ctx:     ctx,
		service: service,

		tx:   tx,
		txID: txID,
		txRow: txmodels.NewTransaction{
			ID:       txID,
			TxStatus: txmodels.TxStatusCreated,
		},

		operations: map[string]*txmodels.NewOperation{},
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

func (f *txFlow) processInputs() ([]txmodels.TrackedOutput, error) {
	outpoints := func(yield func(outpoint bsv.Outpoint) bool) {
		for _, input := range f.tx.Inputs {
			yield(bsv.Outpoint{
				TxID: input.SourceTXID.String(),
				Vout: input.SourceTxOutIndex,
			})
		}
	}

	trackedOutputs, err := f.service.outputs.FindByOutpoints(f.ctx, outpoints)
	if err != nil {
		return nil, txerrors.ErrGettingOutputs.Wrap(err)
	}

	// Check for double-spending
	for _, output := range trackedOutputs {
		if output.IsSpent() && output.SpendingTX != f.txID {
			return nil, txerrors.ErrUTXOSpent
		}
	}

	f.txRow.AddInputs(trackedOutputs...)

	return trackedOutputs, nil
}

func (f *txFlow) operationOfUser(userID string, operationType string, counterparty string) *txmodels.NewOperation {
	if _, ok := f.operations[userID]; !ok {
		f.operations[userID] = &txmodels.NewOperation{
			UserID:       userID,
			Type:         operationType,
			Counterparty: counterparty,

			Transaction: &f.txRow,
			Value:       0,
		}
	}
	return f.operations[userID]
}

func (f *txFlow) addOutputs(outputs ...txmodels.NewOutput) {
	f.txRow.AddOutputs(outputs...)
}

func (f *txFlow) findRelevantP2PKHOutputs() (iter.Seq[txmodels.NewOutput], error) {
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

	rows, err := f.service.addresses.FindByStringAddresses(f.ctx, maps.Keys(relevantOutputs))
	if err != nil {
		return nil, txerrors.ErrGettingAddresses.Wrap(err)
	}

	return func(yield func(output txmodels.NewOutput) bool) {
		for _, row := range rows {
			vout, ok := relevantOutputs[row.Address]
			if !ok {
				f.service.logger.Warn().Str("address", row.Address).Msg("Got not relevant address from database")
				continue
			}
			yield(txmodels.NewOutputForP2PKH(
				bsv.Outpoint{TxID: f.txID, Vout: vout},
				row.UserID,
				bsv.Satoshis(f.tx.Outputs[vout].Satoshis),
				row.CustomInstructions,
			))
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
	return f.service.SaveOperations(f.ctx, maps.Values(f.operations))
}
