package record

import (
	"context"
	"iter"
	"maps"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	database2 "github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/datatypes"
)

type p2pkhOutput struct {
	vout               uint32
	customInstructions datatypes.JSONSlice[bsv.CustomInstruction]
	address            string
	satoshis           bsv.Satoshis
	userID             string
}

type txFlow struct {
	ctx     context.Context
	service *Service

	tx    *trx.Transaction
	txRow *database2.TrackedTransaction
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
		txRow: &database2.TrackedTransaction{
			ID:       txID,
			TxStatus: database2.TxStatusCreated,
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

func (f *txFlow) getFromInputs() ([]*database2.TrackedOutput, error) {
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

	for _, output := range trackedOutputs {
		if output.IsSpent() {
			return nil, txerrors.ErrUTXOSpent
		}
	}

	return trackedOutputs, nil
}

func (f *txFlow) operationOfUser(userID string, operationType string, counterparty string) *operationWrapper {
	if _, ok := f.operations[userID]; !ok {
		f.operations[userID] = &operationWrapper{
			entity: &database2.Operation{
				UserID:       userID,
				Type:         operationType,
				Counterparty: counterparty,

				Transaction: f.txRow,
				Value:       0,
			},
		}
	}
	return f.operations[userID]
}

func (f *txFlow) spendInputs(trackedOutputs []*database2.TrackedOutput) {
	f.txRow.AddInputs(trackedOutputs...)
}

func (f *txFlow) createP2PKHOutput(outputData *p2pkhOutput) {
	f.txRow.CreateP2PKHOutput(&database2.TrackedOutput{
		TxID:     f.txID,
		Vout:     outputData.vout,
		UserID:   outputData.userID,
		Satoshis: outputData.satoshis,
	}, outputData.customInstructions)
}

func (f *txFlow) createDataOutputs(userID string, dataRecords ...*database2.Data) {
	for _, data := range dataRecords {
		f.txRow.CreateDataOutput(data, userID)
	}
}

func (f *txFlow) findRelevantP2PKHOutputs() (iter.Seq[*p2pkhOutput], error) {
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

	return func(yield func(*p2pkhOutput) bool) {
		for _, row := range rows {
			vout, ok := relevantOutputs[row.Address]
			if !ok {
				f.service.logger.Warn().Str("address", row.Address).Msg("Got not relevant address from database")
				continue
			}
			yield(&p2pkhOutput{
				vout:               vout,
				customInstructions: row.CustomInstructions,
				address:            row.Address,
				satoshis:           bsv.Satoshis(f.tx.Outputs[vout].Satoshis),
				userID:             row.UserID,
			})
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
	err := f.service.operations.SaveAll(f.ctx, toOperationEntities(maps.Values(f.operations)))
	if err != nil {
		return txerrors.ErrSavingData.Wrap(err)
	}
	return nil
}
