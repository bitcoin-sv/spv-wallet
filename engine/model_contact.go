package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
)

type Contact struct {
	// Base model
	Model `bson:",inline"`

	// Model specific fields
	ID       string        `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique contact id" bson:"_id"`
	XpubID   string        `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);foreignKey:XpubID;reference:ID;index;comment:This is the related xPub" bson:"xpub_id"`
	FullName string        `json:"full_name" toml:"full_name" yaml:"full_name" gorm:"<-create;comment:This is the contact's full name" bson:"full_name"`
	Paymail  string        `json:"paymail" toml:"paymail" yaml:"paymail" gorm:"<-create;comment:This is the paymail address alias@domain.com" bson:"paymail"`
	PubKey   string        `json:"pub_key" toml:"pub_key" yaml:"pub_key" gorm:"<-:create;index;comment:This is the related public key" bson:"pub_key"`
	Status   ContactStatus `json:"status" toml:"status" yaml:"status" gorm:"<-create;type:varchar(20);default:not authenticated;comment:This is the contact status" bson:"status"`
}

func newContact(fullName, paymailAddress, senderPubKey string, opts ...ModelOps) (*Contact, error) {
	if fullName == "" {
		return nil, ErrEmptyContactFullName
	}

	err := paymail.ValidatePaymail(paymailAddress)

	if err != nil {
		return nil, err
	}

	if senderPubKey == "" {
		return nil, ErrEmptyContactPubKey
	}

	xPubId := utils.Hash(senderPubKey)

	id := utils.Hash(xPubId + paymailAddress)

	contact := &Contact{ID: id, XpubID: xPubId, Model: *NewBaseModel(ModelContact, opts...), FullName: fullName, Paymail: paymailAddress}

	return contact, nil
}

func getContact(ctx context.Context, fullName, paymailAddress, senderPubKey string, opts ...ModelOps) (*Contact, error) {
	contact, err := newContact(fullName, paymailAddress, senderPubKey, opts...)
	if err != nil {
		return nil, err
	}
	contact.ID = ""
	conditions := map[string]interface{}{
		fullNameField: fullName,
		paymailField:  paymailAddress,
	}

	if err := Get(ctx, contact, conditions, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return contact, nil
}

func (c *Contact) getContactPaymailCapability(ctx context.Context) (*paymail.CapabilitiesPayload, error) {
	address := newPaymail(c.Paymail)

	cs := c.Client().Cachestore()
	pc := c.Client().PaymailClient()

	capabilities, err := getCapabilities(ctx, cs, pc, address.Domain)

	if err != nil {
		if errors.Is(err, cachestore.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return capabilities, nil
}

func (c *Contact) getPubKeyFromPki(pkiUrl string) (string, error) {
	if pkiUrl == "" {
		return "", errors.New("pkiUrl should not be empty")
	}
	alias, domain, _ := paymail.SanitizePaymail(c.Paymail)
	pc := c.Client().PaymailClient()

	pkiResponse, err := pc.GetPKI(pkiUrl, alias, domain)

	if err != nil {
		return "", fmt.Errorf("error getting public key from PKI: %w", err)
	}
	return pkiResponse.PubKey, nil
}

func (c *Contact) GetModelName() string {
	return ModelContact.String()
}

// GetModelTableName returns the model db table name
func (c *Contact) GetModelTableName() string {
	return tableContacts
}

// Save the model
func (c *Contact) Save(ctx context.Context) (err error) {
	return Save(ctx, c)
}

// GetID will get the ID
func (c *Contact) GetID() string {
	return c.ID
}

// BeforeCreating is called before the model is saved to the DB
func (c *Contact) BeforeCreating(_ context.Context) (err error) {
	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("starting: %s BeforeCreate hook...", c.Name())

	if c.ID == "" {
		return ErrMissingContactID
	}

	if len(c.FullName) == 0 {
		return ErrMissingContactFullName
	}

	if len(c.Paymail) == 0 {
		return ErrMissingContactPaymail
	}

	if len(c.PubKey) == 0 {
		return ErrMissingContactXPubKey
	}

	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("end: %s BeforeCreate hook", c.Name())
	return
}

// AfterCreated will fire after the model is created in the Datastore
func (c *Contact) AfterCreated(_ context.Context) error {
	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("end: %s AfterCreated hook", c.Name())

	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("end: %s AfterCreated hook", c.Name())
	return nil
}

// Migrate model specific migration on startup
func (c *Contact) Migrate(client datastore.ClientInterface) error {

	tableName := client.GetTableName(tableContacts)
	if client.Engine() == datastore.MySQL {
		if err := c.migrateMySQL(client, tableName); err != nil {
			return err
		}
	} else if client.Engine() == datastore.PostgreSQL {
		if err := c.migratePostgreSQL(client, tableName); err != nil {
			return err
		}
	}

	return client.IndexMetadata(client.GetTableName(tableContacts), MetadataField)
}

// migratePostgreSQL is specific migration SQL for Postgresql
func (c *Contact) migratePostgreSQL(client datastore.ClientInterface, tableName string) error {
	idxName := "idx_" + tableName + "_contacts"
	tx := client.Execute(`CREATE INDEX IF NOT EXISTS "` + idxName + `" ON "` + tableName + `" ("full_name", "paymail")`)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// migrateMySQL is specific migration SQL for MySQL
func (c *Contact) migrateMySQL(client datastore.ClientInterface, tableName string) error {
	idxName := "idx_" + tableName + "_contacts"
	idxExists, err := client.IndexExists(tableName, idxName)
	if err != nil {
		return err
	}
	if !idxExists {
		tx := client.Execute("CREATE UNIQUE INDEX " + idxName + " ON `" + tableName + "` (full_name, paymail)")
		if tx.Error != nil {
			c.Client().Logger().Error().Msgf("failed creating json index on mysql: %s", tx.Error.Error())
			return nil //nolint:nolintlint,nilerr // error is not needed
		}
	}
	return nil
}
