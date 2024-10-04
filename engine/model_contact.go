package engine

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/google/uuid"
)

// Contact is a model that represents a known contacts of the user and invitations to contact.
type Contact struct {
	// Base model
	Model

	// Model specific fields
	ID          string        `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(36);primaryKey;comment:This is the unique contact id"`
	OwnerXpubID string        `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"column:xpub_id;<-:create;type:char(64);foreignKey:XpubID;reference:ID;index;comment:This is the related xPub"`
	FullName    string        `json:"full_name" toml:"full_name" yaml:"full_name" gorm:"<-create;comment:This is the contact's full name"`
	Paymail     string        `json:"paymail" toml:"paymail" yaml:"paymail" gorm:"<-create;comment:This is the paymail address alias@domain.com"`
	PubKey      string        `json:"pub_key" toml:"pub_key" yaml:"pub_key" gorm:"<-:create;index;comment:This is the related public key"`
	Status      ContactStatus `json:"status" toml:"status" yaml:"status" gorm:"<-create;type:varchar(20);default:not confirmed;comment:This is the contact status"`
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

func getContact(ctx context.Context, paymail, ownerXpubID string, opts ...ModelOps) (*Contact, error) {
	paymail = strings.ToLower(paymail)
	conditions := map[string]interface{}{
		xPubIDField:    ownerXpubID,
		paymailField:   paymail,
		deletedAtField: nil,
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

func getContactByID(ctx context.Context, id string, opts ...ModelOps) (*Contact, error) {
	conditions := map[string]interface{}{
		idField: id,
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

// getContactsByXPubIDCount will get a count of all the contacts
func getContactsByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
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

	count, err := getModelCount(
		ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		Contact{}, dbConditions, defaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (m *Contact) validate() error {
	if m.ID == "" {
		return spverrors.ErrMissingContactID
	}

	if m.FullName == "" {
		return spverrors.ErrMissingContactFullName
	}

	if err := paymail.ValidatePaymail(m.Paymail); err != nil {
		return spverrors.ErrInvalidContactPaymail
	}

	if m.PubKey == "" {
		return spverrors.ErrMissingContactXPubKey
	}

	if m.Status == "" {
		return spverrors.ErrMissingContactStatus
	}

	if m.OwnerXpubID == "" {
		return spverrors.ErrMissingContactOwnerXPubID
	}

	return nil
}

func getContacts(ctx context.Context, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Contact, error) {
	if conditions == nil {
		conditions = make(map[string]interface{})
	}

	contacts := make([]*Contact, 0)
	if err := getModelsByConditions(ctx, ModelContact, &contacts, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return contacts, nil
}

func getContactsByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Contact, error) {
	if conditions == nil {
		conditions = make(map[string]interface{})
	}
	conditions[xPubIDField] = xPubID

	contacts := make([]*Contact, 0)
	if err := getModelsByConditions(ctx, ModelContact, &contacts, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Accept marks the contact invitation as accepted, what means that the contact invitation is treated as normal contact.
func (m *Contact) Accept() error {
	if m.Status != ContactAwaitAccept {
		return spverrors.Newf("cannot accept contact. Reason: status: %s, expected: %s", m.Status, ContactAwaitAccept)
	}

	m.Status = ContactNotConfirmed
	return nil
}

// Reject marks the contact invitation as rejected
func (m *Contact) Reject() error {
	if m.Status != ContactAwaitAccept {
		return spverrors.Newf("cannot reject contact. Reason: status: %s, expected: %s", m.Status, ContactAwaitAccept)
	}

	m.DeletedAt.Valid = true
	m.DeletedAt.Time = time.Now()
	m.Status = ContactRejected
	return nil
}

// Confirm marks the contact as confirmed
func (m *Contact) Confirm() error {
	if m.Status != ContactNotConfirmed {
		return spverrors.Newf("cannot confirm contact. Reason: status: %s, expected: %s", m.Status, ContactNotConfirmed)
	}

	m.Status = ContactConfirmed
	return nil
}

// Unconfirm marks the contact as unconfirmed
func (m *Contact) Unconfirm() error {
	if m.Status != ContactConfirmed {
		return spverrors.Newf("cannot unconfirm contact. Reason: status: %s, expected: %s", m.Status, ContactNotConfirmed)
	}

	m.Status = ContactNotConfirmed
	return nil
}

// Delete marks the contact as deleted
func (m *Contact) Delete() {
	m.DeletedAt.Valid = true
	m.DeletedAt.Time = time.Now()
}

// UpdatePubKey updates the contact's public key
func (m *Contact) UpdatePubKey(pk string) (updated bool) {
	if m.PubKey != pk {
		m.PubKey = pk

		if m.Status == ContactConfirmed {
			m.Status = ContactNotConfirmed
		}

		updated = true
		return
	}

	updated = false
	return
}

// GetModelName returns name of the model
func (m *Contact) GetModelName() string {
	return ModelContact.String()
}

// GetModelTableName returns the model db table name
func (m *Contact) GetModelTableName() string {
	return tableContacts
}

// Save the model
func (m *Contact) Save(ctx context.Context) (err error) {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *Contact) GetID() string {
	return m.ID
}

// BeforeCreating is called before the model is saved to the DB
func (m *Contact) BeforeCreating(_ context.Context) (err error) {
	m.Client().Logger().Debug().
		Str("contactID", m.ID).
		Msgf("starting: %s BeforeCreate hook...", m.Name())

	if err = m.validate(); err != nil {
		return
	}

	m.Client().Logger().Debug().
		Str("contactID", m.ID).
		Msgf("end: %s BeforeCreate hook", m.Name())
	return
}

// BeforeUpdating is called before the model is updated in the DB
func (m *Contact) BeforeUpdating(_ context.Context) (err error) {
	m.Client().Logger().Debug().
		Str("contactID", m.ID).
		Msgf("starting: %s BeforeUpdate hook...", m.Name())

	if err = m.validate(); err != nil {
		return
	}

	m.Client().Logger().Debug().
		Str("contactID", m.ID).
		Msgf("end: %s BeforeUpdate hook", m.Name())
	return
}

// Migrate model specific migration on startup
func (m *Contact) Migrate(client datastore.ClientInterface) error {
	tableName := client.GetTableName(tableContacts)
	if err := m.migratePostgreSQL(client, tableName); err != nil {
		return err
	}

	err := client.IndexMetadata(client.GetTableName(tableContacts), MetadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}

// migratePostgreSQL is specific migration SQL for Postgresql
func (m *Contact) migratePostgreSQL(client datastore.ClientInterface, tableName string) error {
	idxName := "idx_" + tableName + "_contacts"
	tx := client.Execute(`CREATE INDEX IF NOT EXISTS "` + idxName + `" ON "` + tableName + `" ("full_name", "paymail")`)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
