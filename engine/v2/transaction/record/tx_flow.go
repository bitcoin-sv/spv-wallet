package record

import (
	"context"
	"iter"
	"maps"

	"github.com/bitcoin-sv/go-sdk/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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

func newTxFlow(ctx context.Context, service *Service, tx *trx.Transaction) (*txFlow, error) {
	txID := tx.TxID().String()
	f := &txFlow{
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

	if err := f.setHex(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *txFlow) setHex() error {
	sourceTXIDs := make([]string, 0, len(f.tx.Inputs))
	for _, input := range f.tx.Inputs {
		sourceTXIDs = append(sourceTXIDs, input.SourceTXID.String())
	}

	foundAll, err := f.service.transactions.HasTransactionInputSources(f.ctx, sourceTXIDs...)
	if err != nil {
		return spverrors.Wrapf(err, "database query failed to check input source transactions for transaction %s", f.txID)
	}

	if foundAll {
		// Optimization: There is no need to serialize the given transaction into BEEFHex format
		// because all its input source ascendants are already stored in the database and serialized as BEEF.
		// This approach avoids unnecessary data redundancy, which impacts overall consistency.
		// All found inputs can be reused when resolving source transaction inputs needed to construct BEEF.
		// (Check the workflow and usage of the BEEF Service in beef_service.go)
		f.txRow.SetRawHex(f.tx.Hex(), sourceTXIDs...)
		return nil
	}

	hex, err := f.tx.BEEFHex()
	if err != nil {
		return spverrors.Wrapf(err, "failed to generate BEEF hex for transaction %s", f.txID)
	}
	f.txRow.SetBEEFHex(hex)

	return nil
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

func (f *txFlow) allP2PKHAddresses() addresses {
	addrs := make(addresses)
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

		addrs.append(address.AddressString, voutU32)
	}
	return addrs
}

func (f *txFlow) createUTXOsForTrackedAddresses(potentialTrackedUniqueAddresses addresses) (iter.Seq[txmodels.NewOutput], error) {
	trackedAddresses, err := f.service.addresses.FindByStringAddresses(f.ctx, maps.Keys(potentialTrackedUniqueAddresses))
	if err != nil {
		return nil, txerrors.ErrGettingAddresses.Wrap(err)
	}

	return func(yield func(output txmodels.NewOutput) bool) {
		for _, tracked := range trackedAddresses {
			addrInfo, ok := potentialTrackedUniqueAddresses[tracked.Address]
			if !ok {
				f.service.logger.Warn().Str("address", tracked.Address).Msg("Got not relevant address from database")
				continue
			}
			for voutsContainingAddress := range addrInfo.vouts {
				yield(txmodels.NewOutputForP2PKH(
					bsv.Outpoint{TxID: f.txID, Vout: voutsContainingAddress},
					tracked.UserID,
					bsv.Satoshis(f.tx.Outputs[voutsContainingAddress].Satoshis),
					tracked.CustomInstructions,
				))
			}
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
	txInfo, err := f.service.broadcaster.Broadcast(f.ctx, f.tx)
	if err != nil {
		return txerrors.ErrTxBroadcast.Wrap(err)
	}

	if txInfo.TXStatus.IsMined() {
		f.txRow.TxStatus = txmodels.TxStatusMined
	} else {
		f.txRow.TxStatus = txmodels.TxStatusBroadcasted
	}

	return nil
}

func (f *txFlow) save() error {
	return f.service.SaveOperations(f.ctx, maps.Values(f.operations))
}
