package engine

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// PaymailAddress is an "external model example" - this model is not part of the standard models loaded and runtime
//
// This model must be included at runtime via WithAutoMigrate() etc...
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type PaymailAddress struct {
	// Base model
	Model

	// Model specific fields
	ID         string `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique paymail record id"`                // Unique identifier
	XpubID     string `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);index;comment:This is the related xPub"`                  // Related xPub ID
	Alias      string `json:"alias" toml:"alias" yaml:"alias" gorm:"<-;type:varchar(64);comment:This is alias@"`                                            // Alias part of the paymail
	Domain     string `json:"domain" toml:"domain" yaml:"domain" gorm:"<-;type:varchar(255);comment:This is @domain.com"`                                   // Domain of the paymail
	PublicName string `json:"public_name" toml:"public_name" yaml:"public_name" gorm:"<-;type:varchar(255);comment:This is public name for public profile"` // Full username
	Avatar     string `json:"avatar" toml:"avatar" yaml:"avatar" gorm:"<-;type:text;comment:This is avatar url"`                                            // This is the url of the user (public profile)

	ExternalXpubKey    string `json:"external_xpub_key" toml:"external_xpub_key" yaml:"external_xpub_key" gorm:"<-:create;type:varchar(512);index;comment:This is full xPub for external use, encryption optional"` // PublicKey hex encoded
	ExternalXpubKeyNum uint32 `json:"external_xpub_num" toml:"external_xpub_num" yaml:"external_xpub_num" gorm:"<-;type:int;default:0;comment:Derivation number used to generate ExternalXpubKey:external_xpub_num"`
	PubKeyNum          uint32 `json:"pubkey_num" toml:"pubkey_num" yaml:"pubkey_num" gorm:"<-;type:int;default:0;comment:Derivation number use to create PKI public key:pubkey_num"`
	XpubDerivationSeq  uint32 `json:"xpub_derivation_seq" toml:"xpub_derivation_seq" yaml:"xpub_derivation_seq" gorm:"<-;type:int;default:0;comment:The index derivation number use to generate new external xpub child keys and rotate PubKey:xpub_derivation_seq"`

	// Private fields
	externalXpubKeyDecrypted string
	externalHdXpub           *compat.ExtendedKey
}

// newPaymail create new paymail model
func newPaymail(paymailAddress string, externalXpubDerivation uint32, opts ...ModelOps) *PaymailAddress {
	// Standardize and sanitize!
	alias, domain, _ := paymail.SanitizePaymail(paymailAddress)
	id, _ := utils.RandomHex(32)
	p := &PaymailAddress{
		Alias:  alias,
		Domain: domain,
		ID:     id,
		Model:  *NewBaseModel(ModelPaymailAddress, opts...),
	}

	// Set the xPub information if found
	if len(p.rawXpubKey) > 0 {
		_ = p.setXPub(externalXpubDerivation)
	}
	return p
}

// getPaymailAddress will get the paymail with the given conditions
func getPaymailAddress(ctx context.Context, address string, opts ...ModelOps) (*PaymailAddress, error) {
	// Get the record
	paymailAddress := &PaymailAddress{Model: *NewBaseModel(ModelPaymailAddress, opts...)}

	alias, domain, _ := paymail.SanitizePaymail(address)
	conditions := map[string]interface{}{
		aliasField:  alias,
		domainField: domain,
	}

	if err := Get(
		ctx, paymailAddress, conditions, false, defaultDatabaseReadTimeout, false,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return paymailAddress, nil
}

// getPaymailAddresses will get all the paymail addresses with the given conditions
func getPaymailAddresses(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*PaymailAddress, error) {
	modelItems := make([]*PaymailAddress, 0)
	if err := getModelsByConditions(ctx, ModelPaymailAddress, &modelItems, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// getPaymailAddressesCount will get all the paymail addresses with the given conditions
func getPaymailAddressesCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelPaymailAddress, PaymailAddress{}, metadata, conditions, opts...)
}

// getPaymailAddressByID will get the paymail with the given ID
func getPaymailAddressByID(ctx context.Context, id string, opts ...ModelOps) (*PaymailAddress, error) {
	// Get the record
	paymailAddress := &PaymailAddress{
		ID:    id,
		Model: *NewBaseModel(ModelPaymailAddress, opts...),
	}
	if err := Get(
		ctx, paymailAddress, nil, false, defaultDatabaseReadTimeout, false,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return paymailAddress, nil
}

// setXPub will set the "ExternalXPubKey" given the raw xPub and xPubID
// encrypted with the given encryption key (if a key is set)
func (m *PaymailAddress) setXPub(externalXpubDerivation uint32) error {
	// Set the ID
	m.XpubID = utils.Hash(m.rawXpubKey)

	// Derive the public key from string
	xPub, err := compat.GetHDKeyFromExtendedPublicKey(m.rawXpubKey)
	if err != nil {
		return spverrors.Wrapf(err, "failed to load xPub key for paymail")
	}

	// Get the external public key
	paymailExternalXpub, err := compat.GetHDKeyByPath(xPub, utils.ChainExternal, externalXpubDerivation)
	if err != nil {
		return spverrors.Wrapf(err, "failed to derive xPub key for paymail")
	}

	m.ExternalXpubKeyNum = externalXpubDerivation

	// Set the decrypted version
	m.externalXpubKeyDecrypted = paymailExternalXpub.String()

	// Encrypt the xPub
	if len(m.encryptionKey) > 0 {
		m.ExternalXpubKey, err = utils.Encrypt(m.encryptionKey, m.externalXpubKeyDecrypted)
	} else {
		m.ExternalXpubKey = m.externalXpubKeyDecrypted
	}

	return spverrors.Wrapf(err, "failed to encrypt external xPub")
}

// GetIdentityXpub will get the identity related to the xPub
func (m *PaymailAddress) GetIdentityXpub() (*compat.ExtendedKey, error) {
	// Get the external xPub (to derive the identity key)
	xPub, err := m.getExternalXpub()
	if err != nil {
		return nil, err
	}

	// Get the last possible key in the external key
	child, err := compat.GetHDKeyChild(
		xPub, uint32(utils.MaxInt32),
	)
	return child, spverrors.Wrapf(err, "failed to generate identity xPub")
}

// getExternalXpub will get the external xPub
func (m *PaymailAddress) getExternalXpub() (*compat.ExtendedKey, error) {
	if m.externalHdXpub != nil {
		return m.externalHdXpub, nil
	}

	// Check if the xPub was encrypted
	if len(m.ExternalXpubKey) != utils.XpubKeyLength {
		var err error
		if m.externalXpubKeyDecrypted, err = utils.Decrypt(
			m.encryptionKey, m.ExternalXpubKey,
		); err != nil {
			return nil, spverrors.Wrapf(err, "failed to decrypt external xPub")
		}
	} else {
		m.externalXpubKeyDecrypted = m.ExternalXpubKey
	}

	// Get the xPub
	xPub, err := compat.GetHDKeyFromExtendedPublicKey(m.externalXpubKeyDecrypted)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get external xPub")
	}

	m.externalHdXpub = xPub
	return m.externalHdXpub, nil
}

// GetPubKey will get the public key for the paymail address.
func (m *PaymailAddress) GetPubKey() (string, error) {
	xPub, err := m.getExternalXpub()
	if err != nil {
		return "", err
	}

	hdPubKey, err := xPub.Child(m.PubKeyNum)
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to derive public key")
	}

	pubKey, err := hdPubKey.ECPubKey()
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to get public key")
	}

	return hex.EncodeToString(pubKey.SerializeCompressed()), nil
}

// GetNextXpub will get the next child xPub for external operations.
func (m *PaymailAddress) GetNextXpub(ctx context.Context) (*compat.ExtendedKey, error) {
	unlock, err := getWaitWriteLockForPaymail(ctx, m.client.Cachestore(), m.ID)
	defer unlock()
	if err != nil {
		return nil, err
	}

	if err = m.incrementExternalXpubDerivationSeq(ctx); err != nil {
		return nil, err
	}

	xPub, err := m.getExternalXpub()
	if err != nil {
		return nil, err
	}

	child, err := xPub.Child(m.XpubDerivationSeq)
	return child, spverrors.Wrapf(err, "failed to get next xPub for external operations")
}

// RotatePubKey will rotate the public key
func (m *PaymailAddress) RotatePubKey(ctx context.Context) error {
	unlock, err := getWaitWriteLockForPaymail(ctx, m.client.Cachestore(), m.ID)
	defer unlock()
	if err != nil {
		return err
	}

	if err = m.incrementExternalXpubDerivationSeq(ctx); err != nil {
		return err
	}

	m.PubKeyNum = m.XpubDerivationSeq
	return m.Save(ctx)
}

// incrementNextNum will atomically update the num of the given chain of the xPub and return it
func (m *PaymailAddress) incrementExternalXpubDerivationSeq(ctx context.Context) error {
	const fieldName = "xpub_derivation_seq"
	var err error
	var newNum int64

	if newNum, err = incrementField(
		ctx, m, fieldName, 1,
	); err != nil {
		return err
	}

	m.XpubDerivationSeq = uint32(newNum)
	return nil
}

// GetModelName returns the model name
func (m *PaymailAddress) GetModelName() string {
	return ModelPaymailAddress.String()
}

// GetModelTableName returns the model db table name
func (m *PaymailAddress) GetModelTableName() string {
	return tablePaymailAddresses
}

// Save the model
func (m *PaymailAddress) Save(ctx context.Context) (err error) {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *PaymailAddress) GetID() string {
	return m.ID
}

// BeforeCreating is called before the model is saved to the DB
func (m *PaymailAddress) BeforeCreating(_ context.Context) (err error) {
	m.Client().Logger().Debug().
		Str("paymailAddressID", m.ID).
		Msgf("starting: %s BeforeCreate hook...", m.Name())

	if m.ID == "" {
		return spverrors.ErrMissingFieldID
	}

	if len(m.Alias) == 0 {
		return spverrors.ErrMissingPaymailAddress
	}

	if len(m.Domain) == 0 {
		return spverrors.ErrMissingPaymailDomain
	}

	if len(m.ExternalXpubKey) == 0 {
		return spverrors.ErrMissingPaymailExternalXPub
	} else if len(m.externalXpubKeyDecrypted) > 0 {
		if _, err = utils.ValidateXPub(m.externalXpubKeyDecrypted); err != nil {
			return
		}
	}

	if len(m.XpubID) == 0 {
		return spverrors.ErrMissingPaymailXPubID
	}

	m.Client().Logger().Debug().
		Str("paymailAddressID", m.ID).
		Msgf("end: %s BeforeCreate hook", m.Name())
	return
}

// AfterCreated will fire after the model is created in the Datastore
func (m *PaymailAddress) AfterCreated(_ context.Context) error {
	m.Client().Logger().Debug().
		Str("paymailAddressID", m.ID).
		Msgf("end: %s AfterCreated hook", m.Name())

	m.Client().Logger().Debug().
		Str("paymailAddressID", m.ID).
		Msgf("end: %s AfterCreated hook", m.Name())
	return nil
}

// Migrate model specific migration on startup
func (m *PaymailAddress) Migrate(client datastore.ClientInterface) error {
	tableName := client.GetTableName(tablePaymailAddresses)
	if err := m.migratePostgreSQL(client, tableName); err != nil {
		return err
	}

	err := client.IndexMetadata(client.GetTableName(tablePaymailAddresses), MetadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}

// migratePostgreSQL is specific migration SQL for Postgresql
func (m *PaymailAddress) migratePostgreSQL(client datastore.ClientInterface, tableName string) error {
	idxName := "idx_" + tableName + "_paymail_address"
	tx := client.Execute(`CREATE INDEX IF NOT EXISTS "` + idxName + `" ON "` + tableName + `" ("alias", "domain")`)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *PaymailAddress) String() string {
	return fmt.Sprintf("%s@%s", m.Alias, m.Domain)
}
