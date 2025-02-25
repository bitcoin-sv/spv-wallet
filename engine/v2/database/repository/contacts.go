package repository

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/dbquery"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/samber/lo"
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
func (r *Contacts) Create(ctx context.Context, newContact *contactsmodels.NewContact) (*contactsmodels.Contact, error) {
	row := database.UserContact{
		UserID:   newContact.UserID,
		FullName: newContact.FullName,
		Paymail:  newContact.NewContactPaymail,
		PubKey:   newContact.NewContactPubKey,
		Status:   newContact.Status,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to create contact")
	}

	return newContactModel(row), nil
}

// Update updates contact in database.
func (r *Contacts) Update(ctx context.Context, contact *contactsmodels.NewContact) (*contactsmodels.Contact, error) {
	row := database.UserContact{
		FullName: contact.FullName,
		PubKey:   contact.NewContactPubKey,
		Status:   contact.Status,
	}
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND paymail = ?", contact.UserID, contact.NewContactPaymail).
		Updates(&row).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to update contact")
	}

	return newContactModel(row), nil
}

// UpdateStatus updates contact status in database.
func (r *Contacts) UpdateStatus(ctx context.Context, userID, paymail, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&database.UserContact{}).
		Where("user_id = ? AND paymail = ?", userID, paymail).
		Update("status", status).Error; err != nil {
		return spverrors.Wrapf(err, "failed to update contact")
	}

	return nil
}

// Delete removes a contact from the database.
func (r *Contacts) Delete(ctx context.Context, userID, paymail string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ? AND paymail = ?", userID, paymail).Delete(&database.UserContact{}).Error; err != nil {
		return spverrors.Wrapf(err, "failed to delete contact")
	}

	return nil
}

// Find retrieves a contact from the database.
func (r *Contacts) Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error) {
	var row = &database.UserContact{}
	if err := r.db.WithContext(ctx).Where("user_id = ? AND paymail = ?", userID, paymail).First(row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

func (r *Contacts) PaginatedForUser(ctx context.Context, userID string, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error) {
	scopes := mapConditionsToScopes(conditions)
	scopes = append(scopes, dbquery.UserID(userID))
	scopes = append(scopes, dbquery.Preload("User"))

	rows, err := dbquery.PaginatedQuery[database.UserContact](
		ctx,
		page,
		r.db,
		scopes...,
	)

	if err != nil {
		return nil, err
	}
	return &models.PagedResult[contactsmodels.Contact]{
		PageDescription: rows.PageDescription,
		Content: lo.Map(rows.Content, func(contact *database.UserContact, _ int) *contactsmodels.Contact {
			return &contactsmodels.Contact{
				ID:        contact.ID,
				UserID:    contact.UserID,
				CreatedAt: contact.CreatedAt,
			}
		}),
	}, nil
}

func newContactModel(row database.UserContact) *contactsmodels.Contact {
	return &contactsmodels.Contact{
		ID:        row.ID,
		UserID:    row.UserID,
		FullName:  row.FullName,
		Paymail:   row.Paymail,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
