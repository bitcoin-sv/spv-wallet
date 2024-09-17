package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// SyncTransaction is an object representing the chain-state sync configuration and results for a given transaction
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type SyncTransaction struct {
	// Base model
	Model `bson:",inline"`

	// Model specific fields
	ID              string      `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique transaction id" bson:"_id"`
	Configuration   SyncConfig  `json:"configuration" toml:"configuration" yaml:"configuration" gorm:"<-;type:text;comment:This is the configuration struct in JSON" bson:"configuration"`
	Results         SyncResults `json:"results" toml:"results" yaml:"results" gorm:"<-;type:text;comment:This is the results struct in JSON" bson:"results"`
	BroadcastStatus SyncStatus  `json:"broadcast_status" toml:"broadcast_status" yaml:"broadcast_status" gorm:"<-;type:varchar(10);index;comment:This is the status of the broadcast" bson:"broadcast_status"`
	P2PStatus       SyncStatus  `json:"p2p_status" toml:"p2p_status" yaml:"p2p_status" gorm:"<-;column:p2p_status;type:varchar(10);index;comment:This is the status of the p2p paymail requests" bson:"p2p_status"`
	SyncStatus      SyncStatus  `json:"sync_status" toml:"sync_status" yaml:"sync_status" gorm:"<-;type:varchar(10);index;comment:This is the status of the on-chain sync" bson:"sync_status"`

	// internal fields
	transaction *Transaction
}

// GetID will get the ID
func (m *SyncTransaction) GetID() string {
	return m.ID
}

// GetModelName will get the name of the current model
func (m *SyncTransaction) GetModelName() string {
	return ModelSyncTransaction.String()
}

// GetModelTableName will get the db table name of the current model
func (m *SyncTransaction) GetModelTableName() string {
	return tableSyncTransactions
}

// Save will save the model into the Datastore
func (m *SyncTransaction) Save(ctx context.Context) error {
	return Save(ctx, m)
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *SyncTransaction) BeforeCreating(_ context.Context) error {
	return nil
}

// AfterCreated will fire after the model is created in the Datastore
func (m *SyncTransaction) AfterCreated(_ context.Context) error {
	return nil
}

// BeforeUpdating will fire before the model is being updated
func (m *SyncTransaction) BeforeUpdating(_ context.Context) error {
	return nil
}

// Migrate model specific migration on startup
func (m *SyncTransaction) Migrate(client datastore.ClientInterface) error {
	err := client.IndexMetadata(client.GetTableName(tableSyncTransactions), metadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}
