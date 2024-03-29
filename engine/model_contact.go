package engine

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/google/uuid"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

type Contact struct {
	// Base model
	Model `bson:",inline"`

	// Model specific fields
	ID          string        `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(36);primaryKey;comment:This is the unique contact id" bson:"_id"`
	OwnerXpubID string        `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"column:xpub_id;<-:create;type:char(64);foreignKey:XpubID;reference:ID;index;comment:This is the related xPub" bson:"xpub_id"`
	FullName    string        `json:"full_name" toml:"full_name" yaml:"full_name" gorm:"<-create;comment:This is the contact's full name" bson:"full_name"`
	Paymail     string        `json:"paymail" toml:"paymail" yaml:"paymail" gorm:"<-create;comment:This is the paymail address alias@domain.com" bson:"paymail"`
	PubKey      string        `json:"pub_key" toml:"pub_key" yaml:"pub_key" gorm:"<-:create;index;comment:This is the related public key" bson:"pub_key"`
	Status      ContactStatus `json:"status" toml:"status" yaml:"status" gorm:"<-create;type:varchar(20);default:not confirmed;comment:This is the contact status" bson:"status"`
}

func newContact(fullName, paymailAddress, pubKey, ownerXpubID string, status ContactStatus, opts ...ModelOps) *Contact {
	contact := Contact{
		Model: *NewBaseModel(ModelContact, opts...),

		ID:          uuid.NewString(),
		OwnerXpubID: ownerXpubID,

		FullName: fullName,
		Paymail:  paymail.SanitizeEmail(paymailAddress),
		PubKey:   pubKey,
		Status:   status,
	}

	return &contact
}

func emptyContact(opts ...ModelOps) *Contact {
	return &Contact{
		Model: *NewBaseModel(ModelContact, opts...),
	}
}

func getContact(ctx context.Context, paymail, ownerXpubID string, opts ...ModelOps) (*Contact, error) {
	conditions := map[string]interface{}{
		xPubIDField:  ownerXpubID,
		paymailField: paymail,
	}

	contact := &Contact{}
	contact.enrich(ModelContact, opts...)

	if err := Get(ctx, contact, conditions, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return contact, nil
}

// GetContactByPaymailAndXPub returns the contact with given paymail and xpub
func getContactByPaymailAndXPubID(ctx context.Context, paymailAddress, xpubID string, opts ...ModelOps) (*Contact, error) {
	contact := emptyContact(opts...)
	contact.Paymail = paymailAddress
	contact.OwnerXpubID = xpubID
	if err := Get(ctx, contact, nil, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}
	return contact, nil
}

func (c *Contact) validate() error {
	if c.ID == "" {
		return ErrMissingContactID
	}

	if c.FullName == "" {
		return ErrMissingContactFullName
	}

	if err := paymail.ValidatePaymail(c.Paymail); err != nil {
		return err
	}

	if c.PubKey == "" {
		return ErrMissingContactXPubKey
	}

	if c.Status == "" {
		return ErrMissingContactStatus
	}

	if c.OwnerXpubID == "" {
		return ErrMissingContactOwnerXPubId
	}

	return nil
}

func getContacts(ctx context.Context, metadata *Metadata, conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Contact, error) {
	contacts := make([]*Contact, 0)

	if err := getModelsByConditions(ctx, ModelContact, &contacts, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return contacts, nil
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

	if err = c.validate(); err != nil {
		return
	}

	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("end: %s BeforeCreate hook", c.Name())
	return
}

func (c *Contact) BeforeUpdating(_ context.Context) (err error) {
	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("starting: %s BeforeUpdate hook...", c.Name())

	if err = c.validate(); err != nil {
		return
	}

	c.Client().Logger().Debug().
		Str("contactID", c.ID).
		Msgf("end: %s BeforeUpdate hook", c.Name())
	return
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
		tx := client.Execute("CREATE INDEX " + idxName + " ON `" + tableName + "` (full_name, paymail)")
		if tx.Error != nil {
			c.Client().Logger().Error().Msgf("failed creating json index on mysql: %s", tx.Error.Error())
			return nil //nolint:nolintlint,nilerr // error is not needed
		}
	}
	return nil
}
