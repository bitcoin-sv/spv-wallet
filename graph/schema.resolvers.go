package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/utils"
)

func (r *mutationResolver) Xpub(ctx context.Context, xpub string, metadata bux.Metadata) (*bux.Xpub, error) {
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

	opts := c.Services.Bux.DefaultModelOptions()
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

func (r *mutationResolver) XpubMetadata(ctx context.Context, metadata bux.Metadata) (*bux.Xpub, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var xPub *bux.Xpub
	xPub, err = c.Services.Bux.UpdateXpubMetadata(ctx, c.XPubID, metadata)
	if err != nil {
		return nil, err
	}

	if !c.Signed || c.XPub == "" {
		xPub.RemovePrivateData()
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

func (r *mutationResolver) AccessKey(ctx context.Context, metadata bux.Metadata) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new accessKey
	var accessKey *bux.AccessKey
	if accessKey, err = c.Services.Bux.NewAccessKey(
		ctx,
		c.XPub,
		bux.WithMetadatas(metadata),
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

func (r *mutationResolver) AccessKeyRevoke(ctx context.Context, id *string) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	// Revoke an accessKey
	var accessKey *bux.AccessKey
	if accessKey, err = c.Services.Bux.RevokeAccessKey(
		ctx,
		c.XPub,
		*id,
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

func (r *mutationResolver) Transaction(ctx context.Context, hex string, draftID *string, metadata bux.Metadata) (*bux.Transaction, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()
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

func (r *mutationResolver) TransactionMetadata(ctx context.Context, id string, metadata bux.Metadata) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx *bux.Transaction
	tx, err = c.Services.Bux.UpdateTransactionMetadata(ctx, c.XPubID, id, metadata)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return bux.DisplayModels(tx).(*bux.Transaction), nil
}

func (r *mutationResolver) NewTransaction(ctx context.Context, transactionConfig bux.TransactionConfig, metadata bux.Metadata) (*bux.DraftTransaction, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	var draftTransaction *bux.DraftTransaction
	draftTransaction, err = c.Services.Bux.NewTransaction(ctx, c.XPub, &transactionConfig, opts...)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(draftTransaction).(*bux.DraftTransaction), nil
}

func (r *mutationResolver) Destination(ctx context.Context, destinationType *string, metadata bux.Metadata) (*bux.Destination, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var useDestinationType string
	if destinationType != nil {
		useDestinationType = *destinationType
	} else {
		useDestinationType = utils.ScriptTypePubKeyHash
	}

	opts := c.Services.Bux.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	var destination *bux.Destination
	destination, err = c.Services.Bux.NewDestination(
		ctx,
		c.XPub,
		utils.ChainExternal,
		useDestinationType,
		true, // monitor this address as it was created by request of a user to share
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

func (r *mutationResolver) DestinationMetadata(ctx context.Context, id *string, address *string, lockingScript *string, metadata bux.Metadata) (*bux.Destination, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var destination *bux.Destination
	if id != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByID(
			ctx,
			c.XPubID,
			*id,
			metadata,
		)
	} else if address != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByAddress(
			ctx,
			c.XPubID,
			*address,
			metadata,
		)
	} else if lockingScript != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByLockingScript(
			ctx,
			c.XPubID,
			*lockingScript,
			metadata,
		)
	}
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

func (r *paymailAddressResolver) PublicName(ctx context.Context, obj *bux.PaymailAddress) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Xpub(ctx context.Context) (*bux.Xpub, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var xPub *bux.Xpub
	xPub, err = c.Services.Bux.GetXpubByID(ctx, c.XPubID)
	if err != nil {
		return nil, err
	}

	if !c.Signed || c.XPub == "" {
		xPub.RemovePrivateData()
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

func (r *queryResolver) AccessKey(ctx context.Context, key string) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var accessKey *bux.AccessKey
	accessKey, err = c.Services.Bux.GetAccessKey(ctx, c.XPubID, key)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

func (r *queryResolver) AccessKeys(ctx context.Context, metadata bux.Metadata) ([]*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var accessKeys []*bux.AccessKey
	accessKeys, err = c.Services.Bux.GetAccessKeys(ctx, c.XPubID, &metadata)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKeys).([]*bux.AccessKey), nil
}

func (r *queryResolver) Transaction(ctx context.Context, id string) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx *bux.Transaction
	tx, err = c.Services.Bux.GetTransaction(ctx, c.XPubID, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return bux.DisplayModels(tx).(*bux.Transaction), nil
}

func (r *queryResolver) Transactions(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) ([]*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx []*bux.Transaction
	tx, err = c.Services.Bux.GetTransactions(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions))
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(tx).([]*bux.Transaction), nil
}

func (r *queryResolver) Destination(ctx context.Context, id *string, address *string, lockingScript *string) (*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var destination *bux.Destination
	if id != nil {
		destination, err = c.Services.Bux.GetDestinationByID(ctx, c.XPubID, *id)
	} else if address != nil {
		destination, err = c.Services.Bux.GetDestinationByAddress(ctx, c.XPubID, *address)
	} else if lockingScript != nil {
		destination, err = c.Services.Bux.GetDestinationByLockingScript(ctx, c.XPubID, *lockingScript)
	} else {
		return nil, bux.ErrMissingFieldID
	}
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
	destinations, err = c.Services.Bux.GetDestinations(ctx, c.XPubID, &metadata)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destinations).([]*bux.Destination), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// PaymailAddress returns generated.PaymailAddressResolver implementation.
func (r *Resolver) PaymailAddress() generated.PaymailAddressResolver {
	return &paymailAddressResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type paymailAddressResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
