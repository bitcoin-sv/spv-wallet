package engine

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/pkg/errors"
)

// UtxoPointer is the actual pointer (index) for the UTXO
type UtxoPointer struct {
	TransactionID string `json:"transaction_id" toml:"transaction_id" yaml:"transaction_id" gorm:"<-:create;type:char(64);index;comment:This is the id of the related transaction" bson:"transaction_id"`
	OutputIndex   uint32 `json:"output_index" toml:"output_index" yaml:"output_index" gorm:"<-:create;type:uint;comment:This is the index of the output in the transaction" bson:"output_index"`
}

// Utxo is an object representing a BitCoin unspent transaction
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Utxo struct {
	// Base model
	Model `bson:",inline"`

	// Standard utxo model base fields
	UtxoPointer `bson:",inline"`

	// Model specific fields
	ID           string                 `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the sha256 hash of the (<txid>|vout)" bson:"_id"`
	XpubID       string                 `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);index;comment:This is the related xPub" bson:"xpub_id"`
	Satoshis     uint64                 `json:"satoshis" toml:"satoshis" yaml:"satoshis" gorm:"<-:create;type:uint;comment:This is the amount of satoshis in the output" bson:"satoshis"`
	ScriptPubKey string                 `json:"script_pub_key" toml:"script_pub_key" yaml:"script_pub_key" gorm:"<-:create;type:text;comment:This is the script pub key" bson:"script_pub_key"`
	Type         string                 `json:"type" toml:"type" yaml:"type" gorm:"<-:create;type:varchar(32);comment:Type of output" bson:"type"`
	DraftID      customTypes.NullString `json:"draft_id" toml:"draft_id" yaml:"draft_id" gorm:"<-;type:varchar(64);index;comment:Related draft id for reservations" bson:"draft_id,omitempty"`
	ReservedAt   customTypes.NullTime   `json:"reserved_at" toml:"reserved_at" yaml:"reserved_at" gorm:"<-;comment:When it was reserved" bson:"reserved_at,omitempty"`
	SpendingTxID customTypes.NullString `json:"spending_tx_id,omitempty" toml:"spending_tx_id" yaml:"spending_tx_id" gorm:"<-;type:char(64);index;comment:This is tx ID of the spend" bson:"spending_tx_id,omitempty"`

	// Virtual field holding the original transaction the utxo originated from
	// This is needed when signing a new transaction that spends the utxo
	Transaction *Transaction `json:"transaction,omitempty" toml:"-" yaml:"-" gorm:"-" bson:"-"`
}

// newUtxo will start a new utxo model
func newUtxo(xPubID, txID, scriptPubKey string, index uint32, satoshis uint64, opts ...ModelOps) *Utxo {
	return &Utxo{
		UtxoPointer: UtxoPointer{
			OutputIndex:   index,
			TransactionID: txID,
		},
		Model:        *NewBaseModel(ModelUtxo, opts...),
		Satoshis:     satoshis,
		ScriptPubKey: scriptPubKey,
		XpubID:       xPubID,
	}
}

// getSpendableUtxos get all spendable utxos by page / pageSize
func getSpendableUtxos(ctx context.Context, xPubID, utxoType string, queryParams *datastore.QueryParams, //nolint:nolintlint,unparam // this param will be used
	fromUtxos []*UtxoPointer, opts ...ModelOps,
) ([]*Utxo, error) {
	// Construct the conditions and results
	var models []Utxo
	conditions := map[string]interface{}{
		draftIDField:      nil,
		spendingTxIDField: nil,
		typeField:         utxoType,
		xPubIDField:       xPubID,
	}

	if fromUtxos != nil {
		for _, fromUtxo := range fromUtxos {
			utxo, err := getUtxo(ctx, fromUtxo.TransactionID, fromUtxo.OutputIndex, opts...)
			if err != nil {
				return nil, err
			} else if utxo == nil {
				return nil, spverrors.ErrCouldNotFindUtxo
			}
			if utxo.XpubID != xPubID || utxo.SpendingTxID.Valid {
				return nil, spverrors.ErrUtxoAlreadySpent
			}
			models = append(models, *utxo)
		}
	} else {
		// Get the records
		if err := getModels(
			ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
			&models, conditions, queryParams, defaultDatabaseReadTimeout,
		); err != nil {
			if errors.Is(err, datastore.ErrNoResults) {
				return nil, nil
			}
			return nil, err
		}
	}

	// No utxos found?
	if len(models) == 0 {
		return nil, spverrors.ErrMissingUTXOsSpendable
	}

	// Loop and enrich
	utxos := make([]*Utxo, 0)
	for index := range models {
		models[index].enrich(ModelUtxo, opts...)
		utxos = append(utxos, &models[index])
	}

	return utxos, nil
}

// unReserveUtxos remove the reservation on the utxos for the given draft ID
func unReserveUtxos(ctx context.Context, xPubID, draftID string, opts ...ModelOps) error {
	var models []Utxo
	conditions := map[string]interface{}{
		xPubIDField:  xPubID,
		draftIDField: draftID,
	}

	// Get the records
	if err := getModels(
		ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		&models, conditions, nil, defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil
		}
		return err
	}

	// Loop and un-reserve
	for index := range models {
		utxo := models[index]
		utxo.enrich(ModelUtxo, opts...)
		utxo.DraftID.Valid = false
		utxo.ReservedAt.Valid = false
		if err := utxo.Save(ctx); err != nil {
			return err
		}
	}

	return nil
}

// reserveUtxos reserve utxos for the given draft ID and amount
func reserveUtxos(ctx context.Context, xPubID, draftID string,
	satoshis uint64, feePerByte float64, fromUtxos []*UtxoPointer, opts ...ModelOps,
) ([]*Utxo, error) {
	// Create base model
	m := NewBaseModel(ModelNameEmpty, opts...)

	// Create the lock and set the release for after the function completes
	unlock, err := newWaitWriteLock(
		ctx, fmt.Sprintf(lockKeyReserveUtxo, xPubID), m.Client().Cachestore(),
	)
	defer unlock()
	if err != nil {
		return nil, err
	}

	// Get spendable utxos
	utxos := new([]*Utxo)
	feeNeeded := uint64(0)
	reservedSatoshis := uint64(0)

	queryParams := &datastore.QueryParams{}
	if fromUtxos == nil {
		// if we are not getting all utxos, paginate the retrieval
		queryParams.Page = 1
		queryParams.PageSize = m.pageSize
		if queryParams.PageSize == 0 {
			queryParams.PageSize = defaultPageSize
		}
	}

reserveUtxoLoop:
	for {
		var freeUtxos []*Utxo
		if freeUtxos, err = getSpendableUtxos(
			ctx, xPubID, utils.ScriptTypePubKeyHash, queryParams, fromUtxos, opts..., // todo: allow reservation of utxos by a different utxo destination type
		); err != nil {
			return nil, err
		}

		if len(freeUtxos) == 0 {
			break reserveUtxoLoop
		}

		// Set vars
		size := utils.GetInputSizeForType(utils.ScriptTypePubKeyHash)

		// Loop the returned utxos
		for _, utxo := range freeUtxos {

			// Set the values on the UTXO
			utxo.DraftID.Valid = true
			utxo.DraftID.String = draftID
			utxo.ReservedAt.Valid = true
			utxo.ReservedAt.Time = time.Now().UTC()

			// Accumulate the reserved satoshis
			reservedSatoshis += utxo.Satoshis

			// Save the UTXO
			// todo: should occur in 1 DB transaction
			if err = utxo.Save(ctx); err != nil {
				return nil, err
			}

			// Add the utxo to the final slice
			*utxos = append(*utxos, utxo)

			// add fee for this new input
			feeNeeded += uint64(float64(size) * feePerByte)
			if reservedSatoshis >= (satoshis + feeNeeded) {
				break reserveUtxoLoop
			}
		}

		if queryParams.PageSize == 0 {
			// break the loop if we are not paginating
			break reserveUtxoLoop
		}
	}

	if reservedSatoshis < satoshis {
		if err = unReserveUtxos(
			ctx, xPubID, draftID, m.GetOptions(false)...,
		); err != nil {
			return nil, spverrors.ErrNotEnoughUtxos
		}
		return nil, spverrors.ErrNotEnoughUtxos
	}

	// check whether an utxo was used twice, this is not valid
	usedUtxos := make([]string, 0)
	for _, utxo := range *utxos {
		if utils.StringInSlice(utxo.ID, usedUtxos) {
			return nil, spverrors.ErrDuplicateUTXOs
		}
		usedUtxos = append(usedUtxos, utxo.ID)
	}

	return *utxos, nil
}

// newUtxoFromTxID will start a new utxo model
func newUtxoFromTxID(txID string, index uint32, opts ...ModelOps) *Utxo {
	return &Utxo{
		Model: *NewBaseModel(ModelUtxo, opts...),
		UtxoPointer: UtxoPointer{
			OutputIndex:   index,
			TransactionID: txID,
		},
	}
}

// getUtxos will get all the utxos with the given conditions
func getUtxos(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Utxo, error) {
	modelItems := make([]*Utxo, 0)
	if err := getModelsByConditions(ctx, ModelUtxo, &modelItems, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// getAccessKeysCount will get a count of all the utxos with the given conditions
func getUtxosCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelUtxo, Utxo{}, metadata, conditions, opts...)
}

// getTransactionsAggregate will get a count of all transactions per aggregate column with the given conditions
func getUtxosAggregate(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	aggregateColumn string, opts ...ModelOps,
) (map[string]interface{}, error) {
	modelItems := make([]*Utxo, 0)
	results, err := getModelsAggregateByConditions(
		ctx, ModelUtxo, &modelItems, metadata, conditions, aggregateColumn, opts...,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// getUtxosByXpubID will return utxos by a given xPub ID
func getUtxosByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Utxo, error) {
	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = conditions
	}
	dbConditions[xPubIDField] = xPubID

	if metadata != nil {
		dbConditions[metadataField] = metadata
	}

	return getUtxosByConditions(ctx, dbConditions, queryParams, opts...)
}

// getUtxosByDraftID will return the utxos by a given draft id
func getUtxosByDraftID(ctx context.Context, draftID string,
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Utxo, error) {
	conditions := map[string]interface{}{
		draftIDField: draftID,
	}
	return getUtxosByConditions(ctx, conditions, queryParams, opts...)
}

// getUtxosByConditions will get utxos by given conditions
func getUtxosByConditions(ctx context.Context, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Utxo, error) {
	var models []Utxo
	if err := getModels(
		ctx, NewBaseModel(
			ModelNameEmpty, opts...).Client().Datastore(),
		&models, conditions, queryParams, databaseLongReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	// Loop and enrich
	utxos := make([]*Utxo, 0)
	for index := range models {
		models[index].enrich(ModelUtxo, opts...)
		utxos = append(utxos, &models[index])
	}
	return utxos, nil
}

// getUtxo will get the utxo with the given conditions
func getUtxo(ctx context.Context, txID string, index uint32, opts ...ModelOps) (*Utxo, error) {
	// Start the new model
	utxo := newUtxoFromTxID(txID, index, opts...)

	// Create the conditions for searching
	conditions := map[string]interface{}{
		"transaction_id": txID,
		"output_index":   index,
	}

	// Get the records
	if err := Get(ctx, utxo, conditions, true, defaultDatabaseReadTimeout, true); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return utxo, nil
}

// GetModelName will get the name of the current model
func (m *Utxo) GetModelName() string {
	return ModelUtxo.String()
}

// GetModelTableName will get the db table name of the current model
func (m *Utxo) GetModelTableName() string {
	return tableUTXOs
}

// Save will save the model into the Datastore
func (m *Utxo) Save(ctx context.Context) (err error) {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *Utxo) GetID() string {
	if m.ID == "" {
		m.ID = m.GenerateID()
	}
	return m.ID
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *Utxo) BeforeCreating(_ context.Context) error {
	m.Client().Logger().Debug().
		Str("utxoID", m.ID).
		Msgf("starting: %s BeforeCreate hook...", m.Name())

	// Test for required field(s)
	if len(m.ScriptPubKey) == 0 {
		return spverrors.ErrMissingFieldScriptPubKey
	} else if m.Satoshis == 0 {
		return spverrors.ErrMissingFieldSatoshis
	} else if len(m.TransactionID) == 0 {
		return spverrors.ErrMissingFieldTransactionID
	}

	if len(m.XpubID) == 0 {
		return spverrors.ErrMissingFieldXpubID
	}

	// Set the new pointer?
	/*
		if m.parsedUtxo == nil {
			m.parsedUtxo = New(bt.UTXO)
		}

		// Parse the UTXO (tx id)
		if m.parsedUtxo.TxID, err = hex.DecodeString(
			m.TransactionID,
		); err != nil {
			return err
		}

		// Parse the UTXO (locking script)
		if m.parsedUtxo.LockingScript, err = bscript2.NewFromHexString(
			m.ScriptPubKey,
		); err != nil {
			return err
		}
		m.parsedUtxo.Satoshis = m.Satoshis
		m.parsedUtxo.Vout = m.OutputIndex
	*/

	// Set the ID
	m.ID = m.GenerateID()
	m.Type = utils.GetDestinationType(m.ScriptPubKey)

	m.Client().Logger().Debug().
		Str("utxoID", m.ID).
		Msgf("end: %s BeforeCreate hook", m.Name())
	return nil
}

// GenerateID will generate the id of the UTXO record based on the format: <txid>|<output_index>
func (m *Utxo) GenerateID() string {
	return utils.Hash(fmt.Sprintf("%s|%d", m.TransactionID, m.OutputIndex))
}

// migratePostgreSQL is specific migration SQL for Postgresql
func (m *Utxo) migratePostgreSQL(client datastore.ClientInterface, tableName string) error {
	tx := client.Execute(`CREATE INDEX IF NOT EXISTS "idx_utxo_reserved" ON "` + tableName + `" ("xpub_id","type","draft_id","spending_tx_id")`)
	return tx.Error
}

// Migrate model specific migration on startup
func (m *Utxo) Migrate(client datastore.ClientInterface) error {
	tableName := client.GetTableName(tableUTXOs)
	if client.Engine() == datastore.PostgreSQL {
		if err := m.migratePostgreSQL(client, tableName); err != nil {
			return err
		}
	}

	return client.IndexMetadata(client.GetTableName(tableUTXOs), metadataField)
}
