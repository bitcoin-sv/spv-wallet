package repository

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"gorm.io/gorm"
)

// Contacts is a repository for contacts.
type Contacts struct {
	db *gorm.DB
}

// NewContactsRepo creates a new repository for addresses.
func NewContactsRepo(db *gorm.DB) *Contacts {
	return &Contacts{db: db}
}

// Create adds a new contact to the database.
func (r *Contacts) Create(ctx context.Context, newContact *contactsmodels.NewContact) error {
	row := &database.UserContact{
		UserID:   newContact.UserID,
		FullName: newContact.FullName,
		Paymail:  newContact.NewContactPaymail,
		PubKey:   newContact.NewContactPubKey,
		Status:   contactsmodels.ContactNotConfirmed,
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return spverrors.Wrapf(err, "failed to create contact")
	}

	return nil
}

// Update updates contact in database.
func (r *Contacts) Update(ctx context.Context, contact *contactsmodels.NewContact) error {
	modelToUpdate := &database.UserContact{
		FullName: contact.FullName,
		PubKey:   contact.NewContactPubKey,
	}
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND paymail = ?", contact.UserID, contact.NewContactPaymail).
		Updates(modelToUpdate).Error; err != nil {
		return spverrors.Wrapf(err, "failed to update contact")
	}

	return nil
}

// Find retrieves a contact from the database.
func (r *Contacts) Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error) {
	var row = &database.UserContact{}
	if err := r.db.WithContext(ctx).Where("user_id = ? AND paymail = ?", userID, paymail).First(row).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to find contact")
	}

	return &contactsmodels.Contact{
		ID:        row.ID,
		UserID:    row.UserID,
		FullName:  row.FullName,
		Paymail:   row.Paymail,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
