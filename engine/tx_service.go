package engine

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// The transaction is treat as external incoming transaction - transaction without a draft
// Only use this function when you know what you are doing!
func saveRawTransaction(ctx context.Context, c ClientInterface, allowUnknown bool, txHex string, opts ...ModelOps) (*Transaction, error) {
	newOpts := c.DefaultModelOptions(append(opts, New())...)
	tx, err := txFromHex(txHex, newOpts...)
	if err != nil {
		return nil, spverrors.ErrMissingTxHex
	}

	// Create the lock and set the release for after the function completes
	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyRecordTx, tx.GetID()), c.Cachestore(),
	)
	defer unlock()
	if err != nil {
		return nil, err
	}

	if !allowUnknown && !tx.hasOneKnownDestination(ctx, c) {
		return nil, spverrors.ErrNoMatchingOutputs
	}

	if err = tx.processUtxos(ctx); err != nil {
		return nil, err
	}

	if !tx.isMined() {
		sync := newSyncTransaction(
			tx.GetID(),
			c.DefaultSyncConfig(),
			tx.GetOptions(true)...,
		)
		sync.BroadcastStatus = SyncStatusSkipped

		sync.Metadata = tx.Metadata
		tx.syncTransaction = sync
	}

	if err = tx.Save(ctx); err != nil {
		return nil, err
	}

	return tx, nil
}

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
	m.NumberOfInputs = uint32(len(m.parsedTx.Inputs))
	m.NumberOfOutputs = uint32(len(m.parsedTx.Outputs))

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
			hex.EncodeToString(m.TransactionBase.parsedTx.Inputs[index].PreviousTxID()),
			m.TransactionBase.parsedTx.Inputs[index].PreviousTxOutIndex,
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
			m.XpubOutputValue[utxo.XpubID] -= int64(utxo.Satoshis)

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

			// only Save utxos for known destinations
			// todo: optimize this SQL SELECT by requesting all the scripts at once (vs in this loop)
			// todo: how to handle tokens and other non-standard outputs ?
			if destination, err = m.transactionService.getDestinationByLockingScript(
				ctx, lockingScript, opts...,
			); err != nil {
				return
			} else if destination != nil {

				// Add value of output to xPub ID
				if _, ok := m.XpubOutputValue[destination.XpubID]; !ok {
					m.XpubOutputValue[destination.XpubID] = 0
				}
				m.XpubOutputValue[destination.XpubID] += int64(amount)

				utxo, _ := m.client.GetUtxoByTransactionID(ctx, m.ID, uint32(i))
				if utxo == nil {
					utxo = newUtxo(
						destination.XpubID, m.ID, txLockingScript, uint32(i),
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
