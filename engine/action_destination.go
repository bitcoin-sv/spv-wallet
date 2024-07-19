package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// NewDestination will get a new destination for an existing xPub
//
// xPubKey is the raw public xPub
func (c *Client) NewDestination(ctx context.Context, xPubKey string, chain uint32,
	destinationType string, opts ...ModelOps,
) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_destination")

	// Get the xPub (by key - converts to id)
	var xPub *Xpub
	var err error
	if xPub, err = getXpubWithCache(
		ctx, c, xPubKey, "", // Get the xPub by xPubID
		c.DefaultModelOptions()..., // Passing down the Datastore and client information into the model
	); err != nil {
		return nil, err
	} else if xPub == nil {
		return nil, spverrors.ErrCouldNotFindXpub
	}

	// Get/create a new destination
	var destination *Destination
	if destination, err = xPub.getNewDestination(
		ctx, chain, destinationType,
		append(opts, c.DefaultModelOptions()...)..., // Passing down the Datastore and client information into the model
	); err != nil {
		return nil, err
	}

	// Save the destination
	if err = destination.Save(ctx); err != nil {
		return nil, err
	}

	// Return the model
	return destination, nil
}

// NewDestinationForLockingScript will create a new destination based on a locking script
func (c *Client) NewDestinationForLockingScript(ctx context.Context, xPubID, lockingScript string,
	opts ...ModelOps,
) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_destination_for_locking_script")

	// Ensure locking script isn't empty
	if len(lockingScript) == 0 {
		return nil, spverrors.ErrMissingLockingScript
	}

	// Start the new destination - will detect type
	destination := newDestination(
		xPubID, lockingScript,
		append(opts, c.DefaultModelOptions()...)..., // Passing down the Datastore and client information into the model
	)

	if destination.Type == "" {
		return nil, spverrors.ErrUnknownLockingScript
	}

	// Save the destination
	if err := destination.Save(ctx); err != nil {
		return nil, err
	}

	// Return the model
	return destination, nil
}

// GetDestinations will get all the destinations from the Datastore
func (c *Client) GetDestinations(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations")

	// Get the destinations
	destinations, err := getDestinations(
		ctx, metadataConditions, conditions, queryParams,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

// GetDestinationsCount will get a count of all the destinations from the Datastore
func (c *Client) GetDestinationsCount(ctx context.Context, metadataConditions *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations_count")

	// Get the destinations count
	count, err := getDestinationsCount(
		ctx, metadataConditions, conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetDestinationsByXpubID will get destinations based on an xPub
//
// metadataConditions are the search criteria used to find destinations
func (c *Client) GetDestinationsByXpubID(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{}, queryParams *datastore.QueryParams,
) ([]*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations")
	// Get the destinations
	destinations, err := getDestinationsByXpubID(
		ctx, xPubID, metadataConditions, conditions, queryParams, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

// GetDestinationsByXpubIDCount will get a count of all destinations based on an xPub
func (c *Client) GetDestinationsByXpubIDCount(ctx context.Context, xPubID string, metadataConditions *Metadata,
	conditions map[string]interface{},
) (int64, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destinations")

	// Get the count
	count, err := getDestinationsCountByXPubID(
		ctx, xPubID, metadataConditions, conditions, c.DefaultModelOptions()...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetDestinationByID will get a destination by id
func (c *Client) GetDestinationByID(ctx context.Context, xPubID, id string) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destination_by_id")

	// Get the destination
	destination, err := getDestinationWithCache(
		ctx, c, id, "", "", c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	// Check that the id matches
	if destination.XpubID != xPubID {
		return nil, spverrors.ErrXpubIDMisMatch
	}

	return destination, nil
}

// GetDestinationByLockingScript will get a destination for a locking script
func (c *Client) GetDestinationByLockingScript(ctx context.Context, xPubID, lockingScript string) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destination_by_locking_script")

	// Get the destination
	destination, err := getDestinationWithCache(
		ctx, c, "", "", lockingScript, c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	// Check that the id matches
	if destination.XpubID != xPubID {
		return nil, spverrors.ErrXpubIDMisMatch
	}

	return destination, nil
}

// GetDestinationByAddress will get a destination for an address
func (c *Client) GetDestinationByAddress(ctx context.Context, xPubID, address string) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "get_destination_by_address")

	// Get the destination
	destination, err := getDestinationWithCache(
		ctx, c, "", address, "", c.DefaultModelOptions()...,
	)
	if err != nil {
		return nil, err
	}

	// Check that the id matches
	if destination.XpubID != xPubID {
		return nil, spverrors.ErrXpubIDMisMatch
	}

	return destination, nil
}

// UpdateDestinationMetadataByID will update the metadata in an existing destination by id
func (c *Client) UpdateDestinationMetadataByID(ctx context.Context, xPubID, id string,
	metadata Metadata,
) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "update_destination_by_id")

	// Get the destination
	destination, err := c.GetDestinationByID(ctx, xPubID, id)
	if err != nil {
		return nil, err
	}

	// Update and save the model
	destination.UpdateMetadata(metadata)
	if err = destination.Save(ctx); err != nil {
		return nil, err
	}

	return destination, nil
}

// UpdateDestinationMetadataByLockingScript will update the metadata in an existing destination by locking script
func (c *Client) UpdateDestinationMetadataByLockingScript(ctx context.Context, xPubID,
	lockingScript string, metadata Metadata,
) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "update_destination_by_locking_script")

	// Get the destination
	destination, err := c.GetDestinationByLockingScript(ctx, xPubID, lockingScript)
	if err != nil {
		return nil, err
	}

	// Update and save the metadata
	destination.UpdateMetadata(metadata)
	if err = destination.Save(ctx); err != nil {
		return nil, err
	}

	return destination, nil
}

// UpdateDestinationMetadataByAddress will update the metadata in an existing destination by address
func (c *Client) UpdateDestinationMetadataByAddress(ctx context.Context, xPubID, address string,
	metadata Metadata,
) (*Destination, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "update_destination_by_address")

	// Get the destination
	destination, err := c.GetDestinationByAddress(ctx, xPubID, address)
	if err != nil {
		return nil, err
	}

	// Update and save the metadata
	destination.UpdateMetadata(metadata)
	if err = destination.Save(ctx); err != nil {
		return nil, err
	}

	return destination, nil
}
