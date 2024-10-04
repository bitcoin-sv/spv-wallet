package engine

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
)

// AccessKey is an object representing an access key model
//
// An AccessKey is a private key with a corresponding public key
// The public key is hashed and saved in this model for retrieval.
// When a request is made with an access key, the public key is sent in the headers, together with
// a signature (like normally done with xPriv signing)
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type AccessKey struct {
	// Base model
	Model

	// Model specific fields
	ID        string               `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique access key id"`
	XpubID    string               `json:"xpub_id" toml:"xpub_id" yaml:"hash" gorm:"<-:create;type:char(64);index;comment:This is the related xPub id"`
	RevokedAt customTypes.NullTime `json:"revoked_at" toml:"revoked_at" yaml:"revoked_at" gorm:"<-;comment:When the key was revoked"`

	// Private fields
	Key string `json:"key" gorm:"-"` // Used on "CREATE", shown to the user "once" only
}

// newAccessKey will start a new model
func newAccessKey(xPubID string, opts ...ModelOps) *AccessKey {
	privateKey, _ := bitcoin.CreatePrivateKey()
	publicKey := hex.EncodeToString(privateKey.PubKey().SerialiseCompressed())
	id := utils.Hash(publicKey)

	return &AccessKey{
		ID:     id,
		Model:  *NewBaseModel(ModelAccessKey, opts...),
		XpubID: xPubID,
		RevokedAt: customTypes.NullTime{NullTime: sql.NullTime{
			Valid: false,
		}},
		Key: hex.EncodeToString(privateKey.Serialise()),
	}
}

// getAccessKey will get the model with a given ID
func getAccessKey(ctx context.Context, id string, opts ...ModelOps) (*AccessKey, error) {
	// Construct an empty tx
	key := &AccessKey{
		ID: id,
	}
	key.enrich(ModelAccessKey, opts...)

	// Get the record
	if err := Get(ctx, key, nil, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}
	return key, nil
}

// getAccessKeys will get all the access keys with the given conditions
func getAccessKeys(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*AccessKey, error) {
	modelItems := make([]*AccessKey, 0)
	if err := getModelsByConditions(ctx, ModelAccessKey, &modelItems, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// getAccessKeysCount will get a count of all the access keys with the given conditions
func getAccessKeysCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelAccessKey, AccessKey{}, metadata, conditions, opts...)
}

// getAccessKeysByXPubID will get all the access keys that match the metadata search
func getAccessKeysByXPubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*AccessKey, error) {
	// Construct an empty model
	var models []AccessKey

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = conditions
	}
	dbConditions[xPubIDField] = xPubID

	if metadata != nil {
		dbConditions[metadataField] = metadata
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
	accessKeys := make([]*AccessKey, 0)
	for index := range models {
		models[index].enrich(ModelDestination, opts...)
		accessKeys = append(accessKeys, &models[index])
	}

	return accessKeys, nil
}

// getAccessKeysByXPubIDCount will get a count of all the access keys that match the metadata search
func getAccessKeysByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = conditions
	}
	dbConditions[xPubIDField] = xPubID

	if metadata != nil {
		dbConditions[metadataField] = metadata
	}

	// Get the records
	count, err := getModelCount(
		ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		AccessKey{}, dbConditions, defaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// GetModelName will get the name of the current model
func (m *AccessKey) GetModelName() string {
	return ModelAccessKey.String()
}

// GetModelTableName will get the db table name of the current model
func (m *AccessKey) GetModelTableName() string {
	return tableAccessKeys
}

// Save will save the model into the Datastore
func (m *AccessKey) Save(ctx context.Context) error {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *AccessKey) GetID() string {
	return m.ID
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *AccessKey) BeforeCreating(_ context.Context) error {
	m.Client().Logger().Debug().
		Str("accessKeyID", m.ID).
		Msgf("starting: %s BeforeCreating hook...", m.Name())

	// Make sure ID is valid
	if len(m.ID) == 0 {
		return spverrors.ErrMissingFieldID
	}

	m.Client().Logger().Debug().
		Str("accessKeyID", m.ID).
		Msgf("end: %s BeforeCreating hook", m.Name())
	return nil
}

// Migrate model specific migration on startup
func (m *AccessKey) Migrate(client datastore.ClientInterface) error {
	err := client.IndexMetadata(client.GetTableName(tableAccessKeys), metadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}
