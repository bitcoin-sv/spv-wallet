package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

// GetDraftTransactionByID will get a draft transaction from the Datastore
func (c *Client) GetDraftTransactionByID(ctx context.Context, id string, opts ...ModelOps) (*DraftTransaction, error) {

	// Get the draft transactions
	draftTransaction, err := getDraftTransactionID(
		ctx, "", id, c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return draftTransaction, nil
}

// GetDraftTransactions will get all the draft transactions from the Datastore
func (c *Client) GetDraftTransactions(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*DraftTransaction, error) {

	// Get the draft transactions
	draftTransactions, err := getDraftTransactions(
		ctx, metadataConditions, conditions, queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return draftTransactions, nil
}

// GetDraftTransactionsCount will get a count of all the draft transactions from the Datastore
func (c *Client) GetDraftTransactionsCount(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {

	// Get the draft transactions count
	count, err := getDraftTransactionsCount(
		ctx, metadataConditions, conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}
