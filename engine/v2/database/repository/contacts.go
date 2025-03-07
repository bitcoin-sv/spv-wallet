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
func (r *Contacts) Create(ctx context.Context, newContact contactsmodels.NewContact) (*contactsmodels.Contact, error) {
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
func (r *Contacts) Update(ctx context.Context, contact contactsmodels.NewContact) (*contactsmodels.Contact, error) {
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

// UpdateByID updates contact full name using its ID.
func (r *Contacts) UpdateByID(ctx context.Context, contactID uint, fullName string) (*contactsmodels.Contact, error) {
	row := database.UserContact{FullName: fullName}
	result := r.db.WithContext(ctx).
		Model(&database.UserContact{}).
		Where("id = ?", contactID).
		Updates(row)

	return r.checkUpdateResultAndReturnContact(ctx, result, contactID)
}

// UpdateStatusByID updates contact status using its ID.
func (r *Contacts) UpdateStatusByID(ctx context.Context, contactID uint, status string) (*contactsmodels.Contact, error) {
	row := database.UserContact{Status: status}
	result := r.db.WithContext(ctx).
		Model(&database.UserContact{}).
		Where("id = ?", contactID).
		Updates(&row)

	return r.checkUpdateResultAndReturnContact(ctx, result, contactID)
}

// Delete removes a contact from the database.
func (r *Contacts) Delete(ctx context.Context, userID, paymail string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ? AND paymail = ?", userID, paymail).Delete(&database.UserContact{}).Error; err != nil {
		return spverrors.Wrapf(err, "failed to delete contact")
	}

	return nil
}

// DeleteByID removes a contact from the database by its ID.
func (r *Contacts) DeleteByID(ctx context.Context, contactID uint) error {
	if err := r.db.WithContext(ctx).Where("id = ?", contactID).Delete(&database.UserContact{}).Error; err != nil {
		return spverrors.Wrapf(err, "failed to delete contact")
	}

	return nil
}

// Find retrieves a contact from the database.
func (r *Contacts) Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error) {
	var row = database.UserContact{}
	if err := r.db.WithContext(ctx).Where("user_id = ? AND paymail = ?", userID, paymail).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, spverrors.Wrapf(err, "failed to find contact")
	}

	return newContactModel(row), nil
}

// FindByID retrieves a contact from the database by id.
func (r *Contacts) FindByID(ctx context.Context, contactID uint) (*contactsmodels.Contact, error) {
	var row = database.UserContact{}
	if err := r.db.WithContext(ctx).Where("id = ?", contactID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, spverrors.Wrapf(err, "failed to find contact")
	}

	return newContactModel(row), nil
}

// PaginatedForUser retrieves contacts for user and the provided paging options and db conditions.
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
			return newContactModel(*contact)
		}),
	}, nil
}

// PaginatedForAdmin retrieves contacts for admin and the provided paging options and db conditions.
func (r *Contacts) PaginatedForAdmin(ctx context.Context, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error) {
	scopes := mapConditionsToScopes(conditions)
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
			return newContactModel(*contact)
		}),
	}, nil
}

func (r *Contacts) checkUpdateResultAndReturnContact(ctx context.Context, result *gorm.DB, contactID uint) (*contactsmodels.Contact, error) {
	if result.Error != nil {
		return nil, spverrors.Wrapf(result.Error, "failed to update contact")
	}

	if result.RowsAffected == 0 {
		return nil, spverrors.ErrUpdateContactStatus
	}

	contact, err := r.FindByID(ctx, contactID)
	if err != nil {
		return nil, spverrors.ErrCannotGetUpdatedContact
	}

	return contact, nil
}

func newContactModel(row database.UserContact) *contactsmodels.Contact {
	contact := &contactsmodels.Contact{
		ID:        row.ID,
		UserID:    row.UserID,
		FullName:  row.FullName,
		Paymail:   row.Paymail,
		Status:    row.Status,
		PubKey:    row.PubKey,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.DeletedAt.Valid {
		contact.DeletedAt = &row.DeletedAt.Time
	}

	return contact
}
