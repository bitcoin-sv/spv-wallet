package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/utils"
)

func (r *mutationResolver) Xpub(ctx context.Context, xpub string, metadata map[string]interface{}) (*bux.Xpub, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var existingXpub *bux.Xpub
	existingXpub, err = c.Services.Bux.GetXpub(ctx, xpub)
	if err != nil && !errors.Is(err, bux.ErrMissingXpub) {
		return nil, err
	}
	if existingXpub != nil {
		return nil, errors.New("xpub already exists")
	}

	opts := make([]bux.ModelOps, 0)
	for key, value := range metadata {
		opts = append(opts, bux.WithMetadata(key, value))
	}

	// Create a new xPub
	var xPub *bux.Xpub
	if xPub, err = c.Services.Bux.NewXpub(
		ctx, xpub, opts...,
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

func (r *mutationResolver) Transaction(ctx context.Context, hex string, draftID *string, metadata map[string]interface{}) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	opts := make([]bux.ModelOps, 0)
	for key, value := range metadata {
		opts = append(opts, bux.WithMetadata(key, value))
	}

	ref := ""
	if draftID != nil {
		ref = *draftID
	}

	var transaction *bux.Transaction
	transaction, err = c.Services.Bux.RecordTransaction(
		ctx, c.XPub, hex, ref, opts...,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicateKey) {
			var txID string
			txID, err = utils.GetTransactionIDFromHex(hex)
			if err != nil {
				return nil, err
			}

			transaction, err = c.Services.Bux.GetTransaction(ctx, c.XPub, txID)
			if err != nil {
				return nil, err
			}

			// record the metadata is being added to the transaction
			if len(metadata) > 0 {
				xPubID := utils.Hash(c.XPub)
				if transaction.XpubMetadata == nil {
					transaction.XpubMetadata = make(bux.XpubMetadata)
				}
				if transaction.XpubMetadata[xPubID] == nil {
					transaction.XpubMetadata[xPubID] = make(bux.Metadata)
				}
				for key, value := range metadata {
					transaction.XpubMetadata[xPubID][key] = value
				}
				err = transaction.Save(ctx)
				if err != nil {
					return nil, err
				}
				// set metadata to the xpub metadata - is removed after Save
				transaction.Metadata = transaction.XpubMetadata[xPubID]
			}

			return transaction, nil
		}
		return nil, err
	}

	return bux.DisplayModels(transaction).(*bux.Transaction), nil
}

func (r *mutationResolver) NewTransaction(ctx context.Context, transactionConfig bux.TransactionConfig, metadata map[string]interface{}) (*bux.DraftTransaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var draftTransaction *bux.DraftTransaction
	draftTransaction, err = c.Services.Bux.NewTransaction(ctx, c.XPub, &transactionConfig, metadata)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(draftTransaction).(*bux.DraftTransaction), nil
}

func (r *mutationResolver) Destination(ctx context.Context, destinationType *string, metadata map[string]interface{}) (*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var useDestinationType string
	if destinationType != nil {
		useDestinationType = *destinationType
	} else {
		useDestinationType = utils.ScriptTypePubKeyHash
	}

	var destination *bux.Destination
	destination, err = c.Services.Bux.NewDestination(
		ctx,
		c.XPub,
		utils.ChainExternal,
		useDestinationType,
		&metadata,
	)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

func (r *queryResolver) Xpub(ctx context.Context) (*bux.Xpub, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var xPub *bux.Xpub
	xPub, err = c.Services.Bux.GetXpub(ctx, c.XPub)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

func (r *queryResolver) Transaction(ctx context.Context, txID string) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx *bux.Transaction
	tx, err = c.Services.Bux.GetTransaction(ctx, c.XPub, txID)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(tx).(*bux.Transaction), nil
}

func (r *queryResolver) Transactions(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) ([]*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx []*bux.Transaction
	tx, err = c.Services.Bux.GetTransactions(ctx, c.XPub, &metadata, ConditionsParseGraphQL(conditions))
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(tx).([]*bux.Transaction), nil
}

func (r *queryResolver) Destination(ctx context.Context, lockingScript string) (*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var destination *bux.Destination
	destination, err = c.Services.Bux.GetDestinationByLockingScript(ctx, c.XPub, lockingScript)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

func (r *queryResolver) Destinations(ctx context.Context, metadata bux.Metadata) ([]*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var destinations []*bux.Destination
	destinations, err = c.Services.Bux.GetDestinations(ctx, c.XPub, &metadata)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destinations).([]*bux.Destination), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
