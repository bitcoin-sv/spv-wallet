package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// GetUtxos will get all the utxos from the Datastore
func (c *Client) GetUtxos(ctx context.Context, metadataConditions *Metadata,
	conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Utxo, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_utxos")

	// Get the utxos
	utxos, err := getUtxos(
		ctx, metadataConditions, conditions, queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	// add the transaction linked to the utxos
	c.enrichUtxoTransactions(ctx, utxos)

	return utxos, nil
}

// GetUtxosCount will get a count of all the utxos from the Datastore
func (c *Client) GetUtxosCount(ctx context.Context, metadataConditions *Metadata,
	conditions *map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_utxos_count")

	// Get the utxos count
	count, err := getUtxosCount(
		ctx, metadataConditions, conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetUtxosByXpubID will get utxos based on an xPub
func (c *Client) GetUtxosByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions *map[string]interface{},
	queryParams *datastore.QueryParams,
) ([]*Utxo, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_utxos")

	// Get the utxos
	utxos, err := getUtxosByXpubID(
		ctx,
		xPubID,
		metadata,
		conditions,
		queryParams,
		c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	// add the transaction linked to the utxos
	c.enrichUtxoTransactions(ctx, utxos)

	return utxos, nil
}

// GetUtxo will get a single utxo based on an xPub, the tx ID and the outputIndex
func (c *Client) GetUtxo(ctx context.Context, xPubKey, txID string, outputIndex uint32) (*Utxo, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_utxo")

	// Get the utxos
	utxo, err := getUtxo(
		ctx, txID, outputIndex, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	} else if utxo == nil {
		return nil, ErrMissingUtxo
	}

	// Check that the id matches
	if utxo.XpubID != utils.Hash(xPubKey) {
		return nil, ErrXpubIDMisMatch
	}

	var tx *Transaction
	tx, err = getTransactionByID(ctx, "", utxo.TransactionID, c.DefaultModelOptions()...)
	if err != nil {
		c.Logger().Error().Str("utxoID", utxo.ID).Msg("failed finding transaction related to utxo")
	} else {
		utxo.Transaction = tx
	}

	return utxo, nil
}

// GetUtxoByTransactionID will get a single utxo based on the tx ID and the outputIndex
func (c *Client) GetUtxoByTransactionID(ctx context.Context, txID string, outputIndex uint32) (*Utxo, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_utxo_by_transaction_id")

	// Get the utxo
	utxo, err := getUtxo(
		ctx, txID, outputIndex, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	} else if utxo == nil {
		return nil, ErrMissingUtxo
	}

	var tx *Transaction
	tx, err = getTransactionByID(ctx, "", utxo.TransactionID, c.DefaultModelOptions()...)
	if err != nil {
		c.Logger().Error().Str("utxoID", utxo.ID).Msg("failed finding transaction related to utxo")
	} else {
		utxo.Transaction = tx
	}

	return utxo, nil
}

// UnReserveUtxos remove the reservation on the utxos for the given draft ID
func (c *Client) UnReserveUtxos(ctx context.Context, xPubID, draftID string) error {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "unreserve_uxtos_by_draft_id")

	return unReserveUtxos(ctx, xPubID, draftID, c.DefaultModelOptions()...)
}

// should this be optional in the results?
func (c *Client) enrichUtxoTransactions(ctx context.Context, utxos []*Utxo) {
	for index, utxo := range utxos {
		tx, err := getTransactionByID(ctx, "", utxo.TransactionID, c.DefaultModelOptions()...)
		if err != nil {
			c.Logger().Error().Str("utxoID", utxo.ID).Msg("failed finding transaction related to utxo")
		} else {
			utxos[index].Transaction = tx
		}
	}
}
