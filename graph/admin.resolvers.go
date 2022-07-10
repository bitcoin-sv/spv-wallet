package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/BuxOrg/bux"
	"github.com/mrz1836/go-datastore"
)

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

func (r *mutationResolver) AdminTransaction(ctx context.Context, hex string) (*bux.Transaction, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()

	var transaction *bux.Transaction
	transaction, err = c.Services.Bux.RecordRawTransaction(
		ctx, hex, opts...,
	)
	if err != nil {
		if !errors.Is(err, datastore.ErrDuplicateKey) {
			return nil, err
		}
	}

	return bux.DisplayModels(transaction).(*bux.Transaction), nil
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
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var accessKeys *bux.AdminStats
	accessKeys, err = c.Services.Bux.GetStats(ctx, c.Services.Bux.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	return accessKeys, nil
}

func (r *queryResolver) AdminAccessKeysList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.AccessKey, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var accessKeys []*bux.AccessKey
	accessKeys, err = c.Services.Bux.GetAccessKeys(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return accessKeys, nil
}

func (r *queryResolver) AdminAccessKeysCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetAccessKeysCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminBlockHeadersList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.BlockHeader, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var blockHeaders []*bux.BlockHeader
	blockHeaders, err = c.Services.Bux.GetBlockHeaders(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return blockHeaders, nil
}

func (r *queryResolver) AdminBlockHeadersCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetBlockHeadersCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminDestinationsList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Destination, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var destinations []*bux.Destination
	destinations, err = c.Services.Bux.GetDestinations(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

func (r *queryResolver) AdminDestinationsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetDestinationsCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminDraftTransactionsList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.DraftTransaction, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var draftTransactions []*bux.DraftTransaction
	draftTransactions, err = c.Services.Bux.GetDraftTransactions(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return draftTransactions, nil
}

func (r *queryResolver) AdminDraftTransactionsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetDraftTransactionsCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
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

func (r *queryResolver) AdminPaymailsList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.PaymailAddress, error) {
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

func (r *queryResolver) AdminPaymailsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
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

func (r *queryResolver) AdminTransactionsList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Transaction, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var transactions []*bux.Transaction
	transactions, err = c.Services.Bux.GetTransactions(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *queryResolver) AdminTransactionsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetTransactionsCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminUtxosList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Utxo, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var utxos []*bux.Utxo
	utxos, err = c.Services.Bux.GetUtxos(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return utxos, nil
}

func (r *queryResolver) AdminUtxosCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetUtxosCount(ctx, &metadata, &conditions)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) AdminXpubsList(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Xpub, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var xpubs []*bux.Xpub
	xpubs, err = c.Services.Bux.GetXPubs(ctx, &metadata, &conditions, params)
	if err != nil {
		return nil, err
	}

	return xpubs, nil
}

func (r *queryResolver) AdminXpubsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
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
