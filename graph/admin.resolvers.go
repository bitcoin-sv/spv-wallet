package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/BuxOrg/bux"
)

func (r *mutationResolver) PaymailCreate(ctx context.Context, address string, publicName *string, avatar *string, metadata bux.Metadata) (*bux.PaymailAddress, error) {
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
	paymailAddress, err = c.Services.Bux.NewPaymailAddress(ctx, c.XPub, address, usePublicName, useAvatar, opts...)
	if err != nil {
		return nil, err
	}

	return paymailAddress, nil
}

func (r *mutationResolver) PaymailDelete(ctx context.Context, address string) (bool, error) {
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
