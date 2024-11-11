package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// Xpub is an object representing an HD-Key or extended public key (xPub for short)
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Xpub struct {
	// Base model
	Model

	// Model specific fields
	ID              string `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the sha256(xpub) hash"`
	CurrentBalance  uint64 `json:"current_balance" toml:"current_balance" yaml:"current_balance" gorm:"<-;comment:The current balance of unspent satoshis"`
	NextInternalNum uint32 `json:"next_internal_num" toml:"next_internal_num" yaml:"next_internal_num" gorm:"<-;type:int;default:0;comment:The index derivation number use to generate NEXT internal xPub (internal xPub are used for change destinations)"`
	NextExternalNum uint32 `json:"next_external_num" toml:"next_external_num" yaml:"next_external_num" gorm:"<-;type:int;default:0;comment:The index derivation number use to generate NEXT external xPub (external xPub are used for address destinations)"`

	destinations []Destination `gorm:"-"` // json:"destinations,omitempty"
}

// newXpub will start a new xPub model
func newXpub(key string, opts ...ModelOps) *Xpub {
	return &Xpub{
		ID:    utils.Hash(key),
		Model: *NewBaseModel(ModelXPub, append(opts, WithXPub(key))...),
	}
}

// newXpubUsingID will start a new xPub model using the xPubID
func newXpubUsingID(xPubID string, opts ...ModelOps) *Xpub {
	return &Xpub{
		ID:    xPubID,
		Model: *NewBaseModel(ModelXPub, opts...),
	}
}

// getXpub will get the xPub with the given conditions
func getXpub(ctx context.Context, key string, opts ...ModelOps) (*Xpub, error) {
	// Get the record
	xPub := newXpub(key, opts...)
	if err := Get(
		ctx, xPub, nil, false, defaultDatabaseReadTimeout, true,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return xPub, nil
}

// getXpubByID will get the xPub with the given conditions
func getXpubByID(ctx context.Context, xPubID string, opts ...ModelOps) (*Xpub, error) {
	// Get the record
	xPub := newXpubUsingID(xPubID, opts...)
	if err := Get(
		ctx, xPub, nil, false, defaultDatabaseReadTimeout, true,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return xPub, nil
}

// getXpubWithCache will try to get from cache first, then datastore
//
// key is the raw xPub key or use xPubID
func getXpubWithCache(ctx context.Context, client ClientInterface,
	key, xPubID string, opts ...ModelOps,
) (*Xpub, error) {
	// Create the cache key
	if len(key) > 0 {
		xPubID = utils.Hash(key)
		opts = append(opts, WithXPub(key)) // Add the xPub option which will set it on the model
	} else if len(xPubID) == 0 {
		return nil, spverrors.ErrMissingFieldXpubID
	}
	cacheKey := fmt.Sprintf(cacheKeyXpubModel, xPubID)

	// Attempt to get from cache
	xPub := new(Xpub)
	found, err := getModelFromCache(
		ctx, client.Cachestore(), cacheKey, xPub,
	)
	if err != nil { // Found a real error
		return nil, err
	} else if found { // Return the cached model
		xPub.enrich(ModelXPub, opts...) // Enrich the model with our parent options
		return xPub, nil
	}

	client.Logger().Info().Str("xpub", xPubID).Msg("xpub not found in cache")

	// Get the xPub
	if xPub, err = getXpubByID(
		ctx, xPubID, opts...,
	); err != nil {
		return nil, err
	} else if xPub == nil {
		return nil, spverrors.ErrCouldNotFindXpub
	}

	// Save to cache
	// todo: run in a go routine
	if err = saveToCache(
		ctx, []string{cacheKey}, xPub, 0,
	); err != nil {
		return nil, err
	}

	// Return the model
	return xPub, nil
}

// getXPubs will get all the xpubs matching the conditions
func getXPubs(ctx context.Context, usingMetadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Xpub, error) {
	modelItems := make([]*Xpub, 0)
	if err := getModelsByConditions(
		ctx, ModelXPub, &modelItems, usingMetadata, conditions, queryParams, opts...,
	); err != nil {
		return nil, err
	}
	return modelItems, nil
}

// getXPubsCount will get a count of the xpubs matching the conditions
func getXPubsCount(ctx context.Context, usingMetadata *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelXPub, Xpub{}, usingMetadata, conditions, opts...)
}

// GetModelName will get the name of the current model
func (m *Xpub) GetModelName() string {
	return ModelXPub.String()
}

// GetModelTableName will get the db table name of the current model
func (m *Xpub) GetModelTableName() string {
	return tableXPubs
}

// Save will save the model into the Datastore
func (m *Xpub) Save(ctx context.Context) error {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *Xpub) GetID() string {
	return m.ID
}

// getNewDestination will get a new destination, adding to the xpub and incrementing num / address
func (m *Xpub) getNewDestination(ctx context.Context, chain uint32, destinationType string,
	opts ...ModelOps,
) (*Destination, error) {
	// Check the type
	// todo: support more types of destinations
	if destinationType != utils.ScriptTypePubKeyHash {
		return nil, spverrors.ErrUnsupportedDestinationType
	}

	// Increment the next num
	num, err := m.getNextDerivationNum(ctx, chain)
	if err != nil {
		return nil, err
	}

	// Create the new address
	var destination *Destination
	if destination, err = newAddress(
		m.rawXpubKey, chain, num, append(opts, New())...,
	); err != nil {
		return nil, err
	}

	// Add the destination to the xPub
	m.destinations = append(m.destinations, *destination)
	return destination, nil
}

// incrementBalance will atomically update the balance of the xPub
func (m *Xpub) incrementBalance(ctx context.Context, balanceIncrement int64) error {
	// Increment the field
	newBalance, err := incrementField(ctx, m, currentBalanceField, balanceIncrement)
	if err != nil {
		return err
	}

	newBalanceU64, err := conv.Int64ToUint64(newBalance)
	if err != nil {
		return spverrors.Wrapf(err, "failed to convert int64 to uint64")
	}
	// Update the field value
	// safe conversion as we have already checked for negative values
	m.CurrentBalance = newBalanceU64

	// Fire the after update
	err = m.AfterUpdated(ctx)
	return err
}

// GetNextInternalDerivationNum will return the next internal derivation number
func (m *Xpub) GetNextInternalDerivationNum(ctx context.Context) (uint32, error) {
	return m.getNextDerivationNum(ctx, utils.ChainInternal)
}

// GetNextExternalDerivationNum will return the next external derivation number
func (m *Xpub) GetNextExternalDerivationNum(ctx context.Context) (uint32, error) {
	return m.getNextDerivationNum(ctx, utils.ChainExternal)
}

func (m *Xpub) getNextDerivationNum(ctx context.Context, chain uint32) (uint32, error) {
	unlock, err := getWaitWriteLockForXpub(ctx, m.client.Cachestore(), m.ID)
	defer unlock()

	if err != nil {
		return 0, err
	}

	derivation, err := m.incrementNextNum(ctx, chain)
	if err != nil {
		return 0, err
	}

	return derivation, nil
}

// incrementNextNum will atomically update the num of the given chain of the xPub and return it
func (m *Xpub) incrementNextNum(ctx context.Context, chain uint32) (uint32, error) {
	var err error
	var newNum int64

	// Choose the field to update
	fieldName := nextExternalNumField
	if chain == utils.ChainInternal {
		fieldName = nextInternalNumField
	}

	// Try to increment the field
	if newNum, err = incrementField(
		ctx, m, fieldName, 1,
	); err != nil {
		return 0, err
	}

	newNumU32, errConversion := conv.Int64ToUint32(newNum)
	if errConversion != nil {
		return 0, spverrors.Wrapf(errConversion, "failed to convert int64 to uint32")
	}

	// Update the model safely as we have already checked for negative values
	if chain == utils.ChainInternal {
		m.NextInternalNum = newNumU32
	} else {
		m.NextExternalNum = newNumU32
	}

	if err = m.AfterUpdated(ctx); err != nil {
		return 0, err
	}

	// Calculate newNumMinusOne
	newNumMinusOne := newNum - 1

	newNumMinusOneU32, errConversion := conv.Int64ToUint32(newNumMinusOne)
	if errConversion != nil {
		return 0, spverrors.Wrapf(errConversion, "failed to convert int64 to uint32")
	}

	// return the previous number, which was next num, safely converted to uint32
	return newNumMinusOneU32, err
}

// ChildModels will get any related sub models
func (m *Xpub) ChildModels() (childModels []ModelInterface) {
	for index := range m.destinations {
		childModels = append(childModels, &m.destinations[index])
	}
	return
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *Xpub) BeforeCreating(_ context.Context) error {
	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("starting: %s BeforeCreating hook...", m.Name())

	// Validate that the xPub key is correct
	if _, err := utils.ValidateXPub(m.rawXpubKey); err != nil {
		return err //nolint:wrapcheck // it is our function returing spverrors
	}

	// Make sure we have an ID
	if len(m.ID) == 0 {
		return spverrors.ErrMissingFieldID
	}

	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("end: %s BeforeCreating hook", m.Name())
	return nil
}

// AfterCreated will fire after the model is created in the Datastore
func (m *Xpub) AfterCreated(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("starting: %s AfterCreated hook...", m.Name())

	// todo: run these in go routines?

	// Store in the cache
	if err := saveToCache(
		ctx, []string{fmt.Sprintf(cacheKeyXpubModel, m.GetID())}, m, 0,
	); err != nil {
		return err
	}

	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("end: %s AfterCreated hook", m.Name())
	return nil
}

// AfterUpdated will fire after a successful update into the Datastore
func (m *Xpub) AfterUpdated(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("starting: %s AfterUpdated hook...", m.Name())

	// Store in the cache
	if err := saveToCache(
		ctx, []string{fmt.Sprintf(cacheKeyXpubModel, m.GetID())}, m, 0,
	); err != nil {
		return err
	}

	m.Client().Logger().Debug().
		Str("xpubID", m.ID).
		Msgf("end: %s AfterUpdated hook", m.Name())
	return nil
}

// PostMigrate is called after the model is migrated
func (m *Xpub) PostMigrate(client datastore.ClientInterface) error {
	err := client.IndexMetadata(client.GetTableName(tableXPubs), metadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}

// RemovePrivateData unset all fields that are sensitive
func (m *Xpub) RemovePrivateData() {
	m.NextExternalNum = 0
	m.NextInternalNum = 0
	m.Metadata = nil
}
