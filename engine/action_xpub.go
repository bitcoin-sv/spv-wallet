package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

// NewXpub will parse the xPub and save it into the Datastore
//
// xPubKey is the raw public xPub
// opts are options and can include "metadata"
func (c *Client) NewXpub(ctx context.Context, xPubKey string, opts ...ModelOps) (*Xpub, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_xpub")

	// Create the model & set the default options (gives options from client->model)
	xPub := newXpub(
		xPubKey, c.DefaultModelOptions(append(opts, New())...)...,
	)

	// Save the model
	if err := xPub.Save(ctx); err != nil {
		return nil, err
	}

	// Return the created model
	return xPub, nil
}

// GetXpub will get an existing xPub from the Datastore
//
// xPubKey is the raw public xPub
func (c *Client) GetXpub(ctx context.Context, xPubKey string) (*Xpub, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_xpub")

	// Attempt to get from cache or datastore
	xPub, err := getXpubWithCache(ctx, c, xPubKey, "", c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	// Return the model
	return xPub, nil
}

// GetXpubByID will get an existing xPub from the Datastore
//
// xPubID is the hash of the xPub
func (c *Client) GetXpubByID(ctx context.Context, xPubID string) (*Xpub, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_xpub_by_id")

	// Attempt to get from cache or datastore
	xPub, err := getXpubWithCache(ctx, c, "", xPubID, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	// Return the model
	return xPub, nil
}

// UpdateXpubMetadata will update the metadata in an existing xPub
//
// xPubID is the hash of the xP
func (c *Client) UpdateXpubMetadata(ctx context.Context, xPubID string, metadata Metadata) (*Xpub, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "update_xpub_by_id")

	// Get the xPub
	xPub, err := c.GetXpubByID(ctx, xPubID)
	if err != nil {
		return nil, err
	}

	// Update the metadata
	xPub.UpdateMetadata(metadata)

	// Save the model
	if err = xPub.Save(ctx); err != nil {
		return nil, err
	}

	// Return the model
	return xPub, nil
}

// GetXPubs gets all xpubs matching the conditions
func (c *Client) GetXPubs(ctx context.Context, metadataConditions *Metadata,
	conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Xpub, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations")

	// Get the count
	xPubs, err := getXPubs(
		ctx, metadataConditions, conditions, queryParams, c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return xPubs, nil
}

// GetXPubsCount gets a count of all xpubs matching the conditions
func (c *Client) GetXPubsCount(ctx context.Context, metadataConditions *Metadata,
	conditions *map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations")

	// Get the count
	count, err := getXPubsCount(
		ctx, metadataConditions, conditions, c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}
