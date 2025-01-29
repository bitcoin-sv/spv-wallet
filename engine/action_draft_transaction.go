package engine

import (
	"context"
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
