package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/mrz1836/go-datastore"
)

// Destination is an object representing a BitCoin destination (address, script, etc)
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Destination struct {
	// Base model
	Model `bson:",inline"`

	// Model specific fields
	ID                           string `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the hash of the locking script" bson:"_id"`
	XpubID                       string `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);index;comment:This is the related xPub" bson:"xpub_id"`
	LockingScript                string `json:"locking_script" toml:"locking_script" yaml:"locking_script" gorm:"<-:create;type:text;comment:This is Bitcoin output script in hex" bson:"locking_script"`
	Type                         string `json:"type" toml:"type" yaml:"type" gorm:"<-:create;type:text;comment:Type of output" bson:"type"`
	Chain                        uint32 `json:"chain" toml:"chain" yaml:"chain" gorm:"<-:create;type:int;comment:This is the (chain)/num location of the address related to the xPub" bson:"chain"`
	Num                          uint32 `json:"num" toml:"num" yaml:"num" gorm:"<-:create;type:int;comment:This is the chain/(num) location of the address related to the xPub" bson:"num"`
	PaymailExternalDerivationNum uint32 `json:"paymail_ext_derivation_num" toml:"paymail_ext_derivation_num" yaml:"paymail_ext_derivation_num" gorm:"<-:create;type:int not null;comment:This is the chain/(num)/(ext_derivation_num) location of the address related to the xPub" bson:"paymail_ext_derivation_num"`
	Address                      string `json:"address" toml:"address" yaml:"address" gorm:"<-:create;type:varchar(35);index;comment:This is the BitCoin address" bson:"address"`
	DraftID                      string `json:"draft_id" toml:"draft_id" yaml:"draft_id" gorm:"<-:create;type:varchar(64);index;comment:This is the related draft id (if internal tx)" bson:"draft_id,omitempty"`
}

// newDestination will start a new Destination model for a locking script
func newDestination(xPubID, lockingScript string, opts ...ModelOps) *Destination {
	// Determine the type if the locking script is provided
	destinationType := utils.ScriptTypeNonStandard
	address := ""
	if len(lockingScript) > 0 {
		destinationType = utils.GetDestinationType(lockingScript)
		address = utils.GetAddressFromScript(lockingScript)
	}

	// Return the model
	return &Destination{
		ID:            utils.Hash(lockingScript),
		LockingScript: lockingScript,
		Model:         *NewBaseModel(ModelDestination, opts...),
		Type:          destinationType,
		XpubID:        xPubID,
		Address:       address,
	}
}

// newAddress will start a new Destination model for a legacy Bitcoin address
func newAddress(rawXpubKey string, chain, num uint32, opts ...ModelOps) (*Destination, error) {
	// Create the model
	destination := &Destination{
		Chain: chain,
		Model: *NewBaseModel(ModelDestination, opts...),
		Num:   num,
	}

	// Set the default address
	err := destination.setAddress(rawXpubKey)
	if err != nil {
		return nil, err
	}

	// Set the locking script
	if destination.LockingScript, err = bitcoin.ScriptFromAddress(
		destination.Address,
	); err != nil {
		return nil, err
	}

	// Determine the type if the locking script is provided
	destination.Type = utils.GetDestinationType(destination.LockingScript)
	destination.ID = utils.Hash(destination.LockingScript)

	// Return the destination (address)
	return destination, nil
}

// getDestinationByID will get the destination by the given id
func getDestinationByID(ctx context.Context, id string, opts ...ModelOps) (*Destination, error) {
	// Construct an empty model
	destination := &Destination{
		ID:    id,
		Model: *NewBaseModel(ModelDestination, opts...),
	}

	// Get the record
	if err := Get(ctx, destination, nil, true, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return destination, nil
}

// getDestinationByAddress will get the destination by the given address
func getDestinationByAddress(ctx context.Context, address string, opts ...ModelOps) (*Destination, error) {
	// Construct an empty model
	destination := &Destination{
		Model: *NewBaseModel(ModelDestination, opts...),
	}
	conditions := map[string]interface{}{
		"address": address,
	}

	// Get the record
	if err := Get(ctx, destination, conditions, true, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return destination, nil
}

// getDestinationByLockingScript will get the destination by the given locking script
func getDestinationByLockingScript(ctx context.Context, lockingScript string, opts ...ModelOps) (*Destination, error) {
	// Construct an empty model
	destination := newDestination("", lockingScript, opts...)

	// Get the record
	if err := Get(ctx, destination, nil, true, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return destination, nil
}

// getDestinations will get all the destinations with the given conditions
func getDestinations(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Destination, error) {
	modelItems := make([]*Destination, 0)
	if err := getModelsByConditions(ctx, ModelDestination, &modelItems, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// getDestinationsCount will get a count of all the destinations with the given conditions
func getDestinationsCount(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelDestination, Destination{}, metadata, conditions, opts...)
}

// getDestinationsByXpubID will get the destination(s) by the given xPubID
func getDestinationsByXpubID(ctx context.Context, xPubID string, usingMetadata *Metadata, conditions *map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Destination, error) {
	// Construct an empty model
	var models []Destination

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = *conditions
	}
	dbConditions[xPubIDField] = xPubID

	if usingMetadata != nil {
		dbConditions[metadataField] = usingMetadata
	}

	// Get the records
	if err := getModels(
		ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		&models, dbConditions, queryParams, defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	// Loop and enrich
	destinations := make([]*Destination, 0)
	for index := range models {
		models[index].enrich(ModelDestination, opts...)
		destinations = append(destinations, &models[index])
	}

	return destinations, nil
}

// getDestinationsCountByXPubID will get a count of the destination(s) by the given xPubID
func getDestinationsCountByXPubID(ctx context.Context, xPubID string, usingMetadata *Metadata,
	conditions *map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = *conditions
	}
	dbConditions[xPubIDField] = xPubID

	if usingMetadata != nil {
		dbConditions[metadataField] = usingMetadata
	}

	// Get the records
	count, err := getModelCount(
		ctx,
		NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		Destination{},
		dbConditions,
		defaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// getXpubWithCache will try to get from cache first, then datastore
//
// key is the raw xPub key or use xPubID
func getDestinationWithCache(ctx context.Context, client ClientInterface,
	id, address, lockingScript string, opts ...ModelOps,
) (*Destination, error) {
	// Create the cache key
	var cacheKey string
	if len(id) > 0 {
		cacheKey = fmt.Sprintf(cacheKeyDestinationModel, id)
	} else if len(address) > 0 {
		cacheKey = fmt.Sprintf(cacheKeyDestinationModelByAddress, address)
	} else if len(lockingScript) > 0 {
		cacheKey = fmt.Sprintf(cacheKeyDestinationModelByLockingScript, lockingScript)
	}
	if len(cacheKey) == 0 {
		return nil, ErrMissingFieldID
	}

	// Attempt to get from cache
	destination := new(Destination)
	found, err := getModelFromCache(
		ctx, client.Cachestore(), cacheKey, destination,
	)
	if err != nil { // Found a real error
		return nil, err
	} else if found { // Return the cached model
		destination.enrich(ModelDestination, opts...) // Enrich the model with our parent options
		return destination, nil
	}

	// Get via ID, address or locking script
	if len(id) > 0 {
		destination, err = getDestinationByID(
			ctx, id, opts...,
		)
	} else if len(address) > 0 {
		destination, err = getDestinationByAddress(
			ctx, address, opts...,
		)
	} else if len(lockingScript) > 0 {
		destination, err = getDestinationByLockingScript(
			ctx, lockingScript, opts...,
		)
	}

	// Check for errors and if the destination is returned
	if err != nil {
		return nil, err
	} else if destination == nil {
		return nil, ErrMissingDestination
	}

	// Save to cache
	// todo: run in a go routine
	if err = saveToCache(
		ctx, []string{
			fmt.Sprintf(cacheKeyDestinationModel, destination.GetID()),
			fmt.Sprintf(cacheKeyDestinationModelByAddress, destination.Address),
			fmt.Sprintf(cacheKeyDestinationModelByLockingScript, destination.LockingScript),
		}, destination, 0,
	); err != nil {
		return nil, err
	}

	// Return the model
	return destination, nil
}

// GetModelName will get the name of the current model
func (m *Destination) GetModelName() string {
	return ModelDestination.String()
}

// GetModelTableName will get the db table name of the current model
func (m *Destination) GetModelTableName() string {
	return tableDestinations
}

// Save will save the model into the Datastore
func (m *Destination) Save(ctx context.Context) (err error) {
	return Save(ctx, m)
}

// GetID will get the model ID
func (m *Destination) GetID() string {
	return m.ID
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *Destination) BeforeCreating(_ context.Context) error {
	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("starting: %s BeforeCreating hook...", m.Name())

	// Set the ID and Type (from LockingScript) (if not set)
	if len(m.LockingScript) > 0 && (len(m.ID) == 0 || len(m.Type) == 0) {
		m.ID = utils.Hash(m.LockingScript)
		m.Type = utils.GetDestinationType(m.LockingScript)
	}

	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("end: %s BeforeCreating hook", m.Name())

	return nil
}

// AfterCreated will fire after the model is created in the Datastore
func (m *Destination) AfterCreated(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("starting: %s AfterCreated hook...", m.Name())

	err := m.client.Cluster().Publish(cluster.DestinationNew, m.LockingScript)
	if err != nil {
		return err
	}

	// Store in the cache
	if err = saveToCache(
		ctx, []string{
			fmt.Sprintf(cacheKeyDestinationModel, m.GetID()),
			fmt.Sprintf(cacheKeyDestinationModelByAddress, m.Address),
			fmt.Sprintf(cacheKeyDestinationModelByLockingScript, m.LockingScript),
		}, m, 0,
	); err != nil {
		return err
	}

	notify(notifications.EventTypeCreate, m)

	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("end: %s AfterCreated hook", m.Name())
	return nil
}

// setAddress will derive and set the address based on the chain (internal vs external)
func (m *Destination) setAddress(rawXpubKey string) error {
	// Check the xPub
	hdKey, err := utils.ValidateXPub(rawXpubKey)
	if err != nil {
		return err
	}

	// Set the ID
	m.XpubID = utils.Hash(rawXpubKey)

	// Derive the address to ensure it is correct
	if m.Address, err = utils.DeriveAddress(
		hdKey, m.Chain, m.Num,
	); err != nil {
		return err
	}

	return nil
}

// Migrate model specific migration on startup
func (m *Destination) Migrate(client datastore.ClientInterface) error {
	return client.IndexMetadata(client.GetTableName(tableDestinations), metadataField)
}

// AfterUpdated will fire after the model is updated in the Datastore
func (m *Destination) AfterUpdated(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("starting: %s AfterUpdated hook...", m.Name())

	// Store in the cache
	if err := saveToCache(
		ctx, []string{
			fmt.Sprintf(cacheKeyDestinationModel, m.GetID()),
			fmt.Sprintf(cacheKeyDestinationModelByAddress, m.Address),
			fmt.Sprintf(cacheKeyDestinationModelByLockingScript, m.LockingScript),
		}, m, 0,
	); err != nil {
		return err
	}

	notify(notifications.EventTypeUpdate, m)

	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("end: %s AfterUpdated hook", m.Name())
	return nil
}

// AfterDeleted will fire after the model is deleted in the Datastore
func (m *Destination) AfterDeleted(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("starting: %s AfterDeleted hook...", m.Name())

	// Only if we have a client, remove all keys
	if m.Client() != nil {
		keys := map[string]string{
			cacheKeyDestinationModel:                m.GetID(),
			cacheKeyDestinationModelByAddress:       m.Address,
			cacheKeyDestinationModelByLockingScript: m.LockingScript,
		}

		for key, val := range keys {
			if err := m.Client().Cachestore().Delete(
				ctx, fmt.Sprintf(key, val),
			); err != nil {
				return err
			}
		}
	}

	notify(notifications.EventTypeDelete, m)

	m.Client().Logger().Debug().
		Str("destinationID", m.ID).
		Msgf("end: %s AfterDeleted hook", m.Name())
	return nil
}
