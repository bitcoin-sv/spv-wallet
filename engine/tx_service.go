package engine

import (
	"context"
	"fmt"
	"math"

	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// processUtxos will process the inputs and outputs for UTXOs
func (m *Transaction) processUtxos(ctx context.Context) error {
	// Input should be processed only for outgoing transactions
	if m.draftTransaction != nil {
		if err := m._processInputs(ctx); err != nil {
			return err
		}
	}

	if err := m._processOutputs(ctx); err != nil {
		return err
	}

	m.TotalValue, m.Fee = m.getValues()
	inputLen, err := conv.IntToUint32(len(m.parsedTx.Inputs))
	if err != nil {
		return spverrors.Wrapf(err, "failed to convert int to uint32")
	}
	m.NumberOfInputs = inputLen

	outputLen, err := conv.IntToUint32(len(m.parsedTx.Outputs))
	if err != nil {
		return spverrors.Wrapf(err, "failed to convert int to uint32")
	}
	m.NumberOfOutputs = outputLen

	return nil
}

// processTxInputs will process the transaction inputs
func (m *Transaction) _processInputs(ctx context.Context) (err error) {
	// Pre-build the options
	opts := m.GetOptions(false)
	client := m.Client()

	var utxo *Utxo

	// check whether we are spending an internal utxo
	for index := range m.TransactionBase.parsedTx.Inputs {
		// todo: optimize this SQL SELECT to get all utxos in one query?
		if utxo, err = m.transactionService.getUtxo(ctx,
			m.TransactionBase.parsedTx.Inputs[index].SourceTXID.String(),
			m.TransactionBase.parsedTx.Inputs[index].SourceTxOutIndex,
			opts...,
		); err != nil {
			return
		} else if utxo != nil { // Found a UTXO record

			// Is Spent?
			if len(utxo.SpendingTxID.String) > 0 {
				return spverrors.ErrUtxoAlreadySpent
			}

			// Only if IUC is enabled (or client is nil which means its enabled by default)
			if client == nil || client.IsIUCEnabled() {

				// check whether the utxo is spent
				isReserved := len(utxo.DraftID.String) > 0
				matchesDraft := m.draftTransaction != nil && utxo.DraftID.String == m.draftTransaction.ID

				// Check whether the spending transaction was reserved by the draft transaction (in the utxo)
				if !isReserved {
					return spverrors.ErrUtxoNotReserved
				}
				if !matchesDraft {
					return spverrors.ErrDraftIDMismatch
				}
			}

			// Update the output value
			if _, ok := m.XpubOutputValue[utxo.XpubID]; !ok {
				m.XpubOutputValue[utxo.XpubID] = 0
			}
			// Check if utxo.Satoshis exceeds int64 range before conversion
			if utxo.Satoshis > math.MaxInt64 {
				return fmt.Errorf("utxo.Satoshis exceeds the maximum value for int64: %d", utxo.Satoshis)
			}
			satoshis, err := conv.Uint64ToInt64(utxo.Satoshis)
			if err != nil {
				return spverrors.Wrapf(err, "failed to convert uint64 to int64")
			}
			m.XpubOutputValue[utxo.XpubID] -= satoshis

			// Mark utxo as spent
			utxo.SpendingTxID.Valid = true
			utxo.SpendingTxID.String = m.ID
			m.utxos = append(m.utxos, *utxo)

			// Add the xPub ID
			if !utils.StringInSlice(utxo.XpubID, m.XpubInIDs) {
				m.XpubInIDs = append(m.XpubInIDs, utxo.XpubID)
			}
		}

		// todo: what if the utxo is nil (not found)?
	}

	return
}

// processTxOutputs will process the transaction outputs
func (m *Transaction) _processOutputs(ctx context.Context) (err error) {
	// Pre-build the options
	opts := m.GetOptions(false)
	newOpts := append(opts, New())
	var destination *Destination

	// check all the outputs for a known destination
	numberOfOutputsProcessed := 0
	for i, output := range m.parsedTx.Outputs {
		amount := output.Satoshis

		// only save outputs with a satoshi value attached to it
		if amount > 0 {

			txLockingScript := output.LockingScript.String()
			lockingScript := utils.GetDestinationLockingScript(txLockingScript)
			address := utils.GetAddressFromScript(lockingScript)

			// only Save utxos for known destinations
			// todo: optimize this SQL SELECT by requesting all the scripts at once (vs in this loop)
			// todo: how to handle tokens and other non-standard outputs ?
			if destination, err = m.transactionService.getDestinationByAddress(
				ctx, address, opts...,
			); err != nil {
				return
			} else if destination != nil {
				i32, err := conv.IntToUint32(i)
				if err != nil {
					return spverrors.Wrapf(err, "failed to convert int to uint32")
				}
				outputIndex := i32

				// Add value of output to xPub ID
				if _, ok := m.XpubOutputValue[destination.XpubID]; !ok {
					m.XpubOutputValue[destination.XpubID] = 0
				}
				amountInt64, err := conv.Uint64ToInt64(amount)
				if err != nil {
					return spverrors.Wrapf(err, "failed to convert uint64 to int64")
				}
				m.XpubOutputValue[destination.XpubID] += amountInt64

				utxo, _ := m.client.GetUtxoByTransactionID(ctx, m.ID, outputIndex)
				if utxo == nil {
					utxo = newUtxo(
						destination.XpubID, m.ID, txLockingScript, outputIndex,
						amount, newOpts...,
					)
				}
				// Append the UTXO model
				m.utxos = append(m.utxos, *utxo)

				// Add the xPub ID
				if !utils.StringInSlice(destination.XpubID, m.XpubOutIDs) {
					m.XpubOutIDs = append(m.XpubOutIDs, destination.XpubID)
				}

				numberOfOutputsProcessed++
			}
		}
	}

	return
}
