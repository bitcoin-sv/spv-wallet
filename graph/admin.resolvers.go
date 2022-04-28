package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/BuxOrg/bux/datastore"
)

func (r *adminStatsResolver) ID(ctx context.Context, obj *bux.AdminStats) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AdminPaymailCreate(ctx context.Context, xpub string, address string, publicName *string, avatar *string, metadata bux.Metadata) (*bux.PaymailAddress, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()

	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	usePublicName := ""
	if publicName != nil {
		usePublicName = *publicName
	}
	useAvatar := ""
	if avatar != nil {
		useAvatar = *avatar
	}

	var paymailAddress *bux.PaymailAddress
	paymailAddress, err = c.Services.Bux.NewPaymailAddress(ctx, xpub, address, usePublicName, useAvatar, opts...)
	if err != nil {
		return nil, err
	}

	return paymailAddress, nil
}

func (r *mutationResolver) AdminPaymailDelete(ctx context.Context, address string) (bool, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return false, err
	}

	opts := c.Services.Bux.DefaultModelOptions()

	// Delete a new paymail address
	err = c.Services.Bux.DeletePaymailAddress(ctx, address, opts...)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) AdminGetStatus(ctx context.Context) (*bool, error) {
	// including admin check
	_, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}

func (r *queryResolver) AdminGetStats(ctx context.Context) (*bux.AdminStats, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AdminPaymailGet(ctx context.Context, address string) (*bux.PaymailAddress, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()

	var paymailAddress *bux.PaymailAddress
	paymailAddress, err = c.Services.Bux.GetPaymailAddress(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return paymailAddress, nil
}

func (r *queryResolver) AdminPaymailList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.PaymailAddress, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var paymailAddresses []*bux.PaymailAddress
	paymailAddresses, err = c.Services.Bux.GetPaymailAddresses(ctx, &metadata, &conditions, nil)
	if err != nil {
		return nil, err
	}

	return paymailAddresses, nil
}

func (r *queryResolver) AdminPaymailCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetPaymailAddressesCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminPaymailGetByXpubID(ctx context.Context, xpubID string) ([]*bux.PaymailAddress, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var paymailAddresses []*bux.PaymailAddress
	paymailAddresses, err = c.Services.Bux.GetPaymailAddressesByXPubID(ctx, xpubID, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return paymailAddresses, nil
}

func (r *queryResolver) AdminXpubList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Xpub, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var xPubs []*bux.Xpub
	xPubs, err = c.Services.Bux.GetXPubs(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return xPubs, nil
}

func (r *queryResolver) AdminXpubCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetXPubsCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// AdminStats returns generated.AdminStatsResolver implementation.
func (r *Resolver) AdminStats() generated.AdminStatsResolver { return &adminStatsResolver{r} }

type adminStatsResolver struct{ *Resolver }
