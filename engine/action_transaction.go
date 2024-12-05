package engine

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"math"
	"time"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// RecordTransaction will parse the outgoing transaction and save it into the Datastore
// xPubKey is the raw public xPub
// txHex is the raw transaction hex
// draftID is the unique draft id from a previously started New() transaction (draft_transaction.ID)
// opts are model options and can include "metadata"
func (c *Client) RecordTransaction(ctx context.Context, xPubKey, txHex, draftID string, opts ...ModelOps) (*Transaction, error) {
	tx, err := trx.NewTransactionFromHex(txHex)
	if err != nil {
		return nil, spverrors.ErrInvalidHex
	}

	rts, err := getOutgoingTxRecordStrategy(xPubKey, tx, draftID)
	if err != nil {
		return nil, err
	}

	return recordTransaction(ctx, c, rts, opts...)
}

// RecordRawTransaction will parse the transaction and save it into the Datastore directly, without any checks or broadcast but SPV Wallet Engine will ask network for information if transaction was mined
// The transaction is treat as external incoming transaction - transaction without a draft
// Only use this function when you know what you are doing!
//
// txHex is the raw transaction hex
// opts are model options and can include "metadata"
func (c *Client) RecordRawTransaction(ctx context.Context, txHex string,
	opts ...ModelOps,
) (*Transaction, error) {

	return saveRawTransaction(ctx, c, true, txHex, opts...)
}

// NewTransaction will create a new draft transaction and return it
//
// ctx is the context
// rawXpubKey is the raw xPub key
// config is the TransactionConfig
// metadata is added to the model
// opts are additional model options to be applied
func (c *Client) NewTransaction(ctx context.Context, rawXpubKey string, config *TransactionConfig,
	opts ...ModelOps,
) (*DraftTransaction, error) {

	// Create the lock and set the release for after the function completes
	unlock, err := getWaitWriteLockForXpub(
		ctx, c.Cachestore(), utils.Hash(rawXpubKey),
	)
	defer unlock()
	if err != nil {
		return nil, err
	}

	// Create the draft tx model
	draftTransaction, err := newDraftTransaction(
		rawXpubKey, config,
		c.DefaultModelOptions(append(opts, New())...)...,
	)
	if err != nil {
		return nil, err
	}

	// Save the model
	if err = draftTransaction.Save(ctx); err != nil {
		return nil, err
	}

	// Return the created model
	return draftTransaction, nil
}

// GetTransaction will get a transaction by its ID from the Datastore
func (c *Client) GetTransaction(ctx context.Context, xPubID, txID string) (*Transaction, error) {
	// Get the transaction by ID
	transaction, err := getTransactionByID(
		ctx, xPubID, txID, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}

	return transaction, nil
}

// GetAdminTransaction will get a transaction by its ID from the Datastore
func (c *Client) GetAdminTransaction(ctx context.Context, txID string) (*Transaction, error) {

	// Get the transaction by ID
	transaction, err := getAdminTransactionByID(
		ctx, txID, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}

	return transaction, nil
}

// GetTransactionsByIDs returns array of transactions by their IDs from the Datastore
func (c *Client) GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error) {

	// Create the conditions
	conditions := generateTxIDFilterConditions(txIDs)

	// Get the transactions by its IDs
	transactions, err := getTransactions(
		ctx, nil, conditions, nil,
		c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionByHex will get a transaction from the Datastore by its full hex string
// uses GetTransaction
func (c *Client) GetTransactionByHex(ctx context.Context, hex string) (*Transaction, error) {
	tx, err := trx.NewTransactionFromHex(hex)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to parse transaction hex: %s", hex)
	}

	return c.GetTransaction(ctx, "", tx.TxID().String())
}

// GetTransactions will get all the transactions from the Datastore
func (c *Client) GetTransactions(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Transaction, error) {

	// Get the transactions
	transactions, err := getTransactions(
		ctx, metadataConditions, conditions, queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionsCount will get a count of all the transactions from the Datastore
func (c *Client) GetTransactionsCount(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {

	// Get the transactions count
	count, err := getTransactionsCount(
		ctx, metadataConditions, conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetTransactionsByXpubID will get all transactions for a given xpub from the Datastore
//
// ctx is the context
// rawXpubKey is the raw xPub key
// metadataConditions is added to the request for searching
// conditions is added the request for searching
func (c *Client) GetTransactionsByXpubID(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams,
) ([]*Transaction, error) {

	// Get the transaction by ID
	transactions, err := getTransactionsByXpubID(
		ctx,
		xPubID,
		metadataConditions,
		conditions,
		queryParams,
		c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionsByXpubIDCount will get the count of all transactions matching the search criteria
func (c *Client) GetTransactionsByXpubIDCount(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{},
) (int64, error) {

	count, err := getTransactionsCountByXpubID(
		ctx, xPubID, metadataConditions, conditions,
		c.DefaultModelOptions()...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateTransactionMetadata will update the metadata in an existing transaction
func (c *Client) UpdateTransactionMetadata(ctx context.Context, xPubID, id string,
	metadata Metadata,
) (*Transaction, error) {

	// Get the transaction
	transaction, err := c.GetTransaction(ctx, xPubID, id)
	if err != nil {
		return nil, err
	}

	// Update the metadata
	if err = transaction.UpdateTransactionMetadata(
		xPubID, metadata,
	); err != nil {
		return nil, err
	}

	// Save the model	// update existing record
	if err = transaction.Save(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

// RevertTransaction will revert a transaction created in the SPV Wallet Engine database, but only if it has not
// yet been synced on-chain and the utxos have not been spent.
// All utxos that are reverted will be marked as deleted (and spent)
func (c *Client) RevertTransaction(ctx context.Context, id string) error {
	// Get the transaction
	transaction, err := c.GetTransaction(ctx, "", id)
	if err != nil {
		return err
	}

	// make sure the transaction is coming from SPV Wallet Engine
	if transaction.DraftID == "" {
		return spverrors.ErrTxRevertEmptyDraftID
	}

	draftTransaction, err := c.GetDraftTransactionByID(ctx, transaction.DraftID, c.DefaultModelOptions()...)
	if err != nil {
		return err
	}
	if draftTransaction == nil {
		return spverrors.ErrTxRevertCouldNotFindDraftTx
	}

	// check that the utxos of this transaction have not been spent
	// this transaction needs to be the tip of the chain
	conditions := map[string]interface{}{
		"transaction_id": transaction.ID,
	}
	var utxos []*Utxo
	if utxos, err = c.GetUtxos(ctx, nil, conditions, nil, c.DefaultModelOptions()...); err != nil {
		return err
	}
	for _, utxo := range utxos {
		if utxo.SpendingTxID.Valid {
			return spverrors.ErrTxRevertUtxoAlreadySpent
		}
	}

	//
	// Revert transaction and all related elements
	//

	// mark output utxos as deleted (no way to delete from SPV Wallet Engine yet)
	for _, utxo := range utxos {
		utxo.enrich(ModelUtxo, c.DefaultModelOptions()...)
		utxo.SpendingTxID.Valid = true
		utxo.SpendingTxID.String = "deleted"
		utxo.DeletedAt.Valid = true
		utxo.DeletedAt.Time = time.Now()
		if err = utxo.Save(ctx); err != nil {
			return err
		}
	}

	// remove output values of transaction from all xpubs
	var xpub *Xpub
	for xpubID, outputValue := range transaction.XpubOutputValue {
		if xpub, err = c.GetXpubByID(ctx, xpubID); err != nil {
			return err
		}
		if outputValue > 0 {
			xpub.CurrentBalance -= uint64(outputValue)
		} else {
			xpub.CurrentBalance += uint64(math.Abs(float64(outputValue)))
		}

		if err = xpub.Save(ctx); err != nil {
			return err
		}
	}

	// set any inputs (spent utxos) used in this transaction back to not spent
	var utxo *Utxo
	for _, input := range draftTransaction.Configuration.Inputs {
		if utxo, err = c.GetUtxoByTransactionID(ctx, input.TransactionID, input.OutputIndex); err != nil {
			return err
		}
		utxo.SpendingTxID.Valid = false
		utxo.SpendingTxID.String = ""
		if err = utxo.Save(ctx); err != nil {
			return err
		}
	}

	// cancel draft tx
	draftTransaction.Status = DraftStatusCanceled
	if err = draftTransaction.Save(ctx); err != nil {
		return err
	}

	// revert transaction
	// this takes the transaction out of any possible list view of the owners of the xpubs,
	// but keeps a record of what went down
	if transaction.Metadata == nil {
		transaction.Metadata = Metadata{}
	}
	transaction.Metadata["XpubInIDs"] = transaction.XpubInIDs
	transaction.Metadata["XpubOutIDs"] = transaction.XpubOutIDs
	transaction.Metadata["XpubOutputValue"] = transaction.XpubOutputValue
	transaction.XpubInIDs = IDs{"reverted"}
	transaction.XpubOutIDs = IDs{"reverted"}
	transaction.XpubOutputValue = XpubOutputValue{"reverted": 0}
	transaction.DeletedAt.Valid = true
	transaction.DeletedAt.Time = time.Now()
	transaction.TxStatus = TxStatusReverted

	err = transaction.Save(ctx) // update existing record

	return err
}

// HandleTxCallback will update the broadcast callback transaction info, like: block height, block hash, status, bump.
func (c *Client) HandleTxCallback(ctx context.Context, callbackResp *chainmodels.TXInfo) error {
	logger := c.options.logger
	bump, err := trx.NewMerklePathFromHex(callbackResp.MerklePath)

	if err != nil {
		logger.Err(err).Msgf("failed to parse merkle path from broadcast callback - tx: %v", callbackResp)
		return spverrors.Wrapf(err, "failed to parse merkle path from broadcast callback - tx: %v", callbackResp)
	}

	txID := callbackResp.TxID

	tx, err := c.GetTransaction(ctx, "", txID)
	if err != nil {
		logger.Warn().Err(err).Msgf("failed to get transaction by id: %v", txID)
		return err
	}

	tx.BlockHash = callbackResp.BlockHash
	blockHeight, err := conv.Int64ToUint64(callbackResp.BlockHeight)
	if err != nil {
		return spverrors.Wrapf(err, "failed to convert block height to uint64 - tx: %v", callbackResp.BlockHeight)
	}
	tx.BlockHeight = blockHeight
	tx.SetBUMP(bump)
	tx.UpdateFromBroadcastStatus(callbackResp.TXStatus)

	if err := tx.Save(ctx); err != nil {
		return spverrors.ErrDuringSaveTx.Wrap(err)
	}

	return nil
}

func (c *Client) GetTransactionData(ctx context.Context, outpoint bsv.Outpoint) ([]byte, error) {
	data, err := c.txDAO.GetData(ctx, outpoint)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get transaction's data")
	}
	return data, nil
}

func generateTxIDFilterConditions(txIDs []string) map[string]interface{} {
	orConditions := make([]map[string]interface{}, len(txIDs))

	for i, txID := range txIDs {
		orConditions[i] = map[string]interface{}{"id": txID}
	}

	return map[string]interface{}{
		"$or": orConditions,
	}
}
