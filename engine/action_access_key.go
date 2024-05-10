package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// NewAccessKey will create a new access key for the given xpub
//
// opts are options and can include "metadata"
func (c *Client) NewAccessKey(ctx context.Context, rawXpubKey string, opts ...ModelOps) (*AccessKey, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_access_key")

	// Validate that the value is an xPub
	_, err := utils.ValidateXPub(rawXpubKey)
	if err != nil {
		return nil, err
	}

	// Get the xPub (by key - converts to id)
	var xPub *Xpub
	if xPub, err = getXpubWithCache(
		ctx, c, rawXpubKey, "", // Pass the context and key everytime (for now)
		c.DefaultModelOptions()..., // Passing down the Datastore and client information into the model
	); err != nil {
		return nil, err
	} else if xPub == nil {
		return nil, ErrMissingXpub
	}

	// Create the model & set the default options (gives options from client->model)
	accessKey := newAccessKey(
		xPub.ID, c.DefaultModelOptions(append(opts, New())...)...,
	)

	// Save the model
	if err = accessKey.Save(ctx); err != nil {
		return nil, err
	}

	// Return the created model
	return accessKey, nil
}

// GetAccessKey will get an existing access key from the Datastore
func (c *Client) GetAccessKey(ctx context.Context, xPubID, id string) (*AccessKey, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_access_key")

	// Get the access key
	accessKey, err := getAccessKey(
		ctx, id,
		c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	} else if accessKey == nil {
		return nil, ErrAccessKeyNotFound
	}

	// make sure this is the correct accessKey
	if accessKey.XpubID != xPubID {
		return nil, utils.ErrXpubNoMatch
	}

	// Return the model
	return accessKey, nil
}

// GetAccessKeys will get all the access keys from the Datastore
func (c *Client) GetAccessKeys(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*AccessKey, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_access_keys")

	// Get the access keys
	accessKeys, err := getAccessKeys(
		ctx, metadataConditions, conditions, queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return accessKeys, nil
}

// GetAccessKeysCount will get a count of all the access keys from the Datastore
func (c *Client) GetAccessKeysCount(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_access_keys_count")

	// Get the access keys count
	count, err := getAccessKeysCount(
		ctx, metadataConditions, conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetAccessKeysByXPubID will get all existing access keys from the Datastore
//
// metadataConditions is the metadata to match to the access keys being returned
func (c *Client) GetAccessKeysByXPubID(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*AccessKey, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_access_keys")

	// Get the access key
	accessKeys, err := getAccessKeysByXPubID(
		ctx,
		xPubID,
		metadataConditions,
		conditions,
		queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	} else if accessKeys == nil {
		return nil, datastore.ErrNoResults
	}

	// Return the models
	return accessKeys, nil
}

// GetAccessKeysByXPubIDCount will get a count of all existing access keys from the Datastore
func (c *Client) GetAccessKeysByXPubIDCount(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_access_keys")

	// Get the access key
	count, err := getAccessKeysByXPubIDCount(
		ctx,
		xPubID,
		metadataConditions,
		conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	// Return the models
	return count, nil
}

// RevokeAccessKey will revoke an access key by its id
//
// opts are options and can include "metadata"
func (c *Client) RevokeAccessKey(ctx context.Context, rawXpubKey, id string, opts ...ModelOps) (*AccessKey, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_access_key")

	// Validate that the value is an xPub
	_, err := utils.ValidateXPub(rawXpubKey)
	if err != nil {
		return nil, err
	}

	// Get the xPub (by key - converts to id)
	var xPub *Xpub
	if xPub, err = getXpubWithCache(
		ctx, c, rawXpubKey, "", // Pass the context and key everytime (for now)
		c.DefaultModelOptions()..., // Passing down the Datastore and client information into the model
	); err != nil {
		return nil, err
	} else if xPub == nil {
		return nil, ErrMissingXpub
	}

	var accessKey *AccessKey
	if accessKey, err = getAccessKey(
		ctx, id, c.DefaultModelOptions(opts...)...,
	); err != nil {
		return nil, err
	}
	if accessKey == nil {
		return nil, ErrMissingAccessKey
	}

	// make sure this is the correct accessKey
	xPubID := utils.Hash(rawXpubKey)
	if accessKey.XpubID != xPubID {
		return nil, utils.ErrXpubNoMatch
	}

	accessKey.RevokedAt.Valid = true
	accessKey.RevokedAt.Time = time.Now()

	// Save the model
	if err = accessKey.Save(ctx); err != nil {
		return nil, err
	}

	// Return the updated model
	return accessKey, nil
}
