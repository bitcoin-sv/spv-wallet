package repository

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	dberrors "github.com/bitcoin-sv/spv-wallet/engine/v2/database/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"gorm.io/gorm"
)

// Paymails is a repository for paymails.
type Paymails struct {
	db *gorm.DB
}

// NewPaymailsRepo creates a new repository for paymails.
func NewPaymailsRepo(db *gorm.DB) *Paymails {
	return &Paymails{db: db}
}

// Create adds a new paymail to the database.
func (p *Paymails) Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error) {
	row := database.Paymail{
		Alias:  newPaymail.Alias,
		Domain: newPaymail.Domain,

		PublicName: newPaymail.PublicName,
		Avatar:     newPaymail.Avatar,

		UserID: newPaymail.UserID,
	}

	if err := p.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, dberrors.QueryFailed.Wrap(err, "failed to create paymail")
	}

	return p.newPaymailModel(row), nil
}

// Find returns a paymail by alias and domain.
func (p *Paymails) Find(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error) {
	var row database.Paymail
	if err := p.db.
		WithContext(ctx).
		Where("alias = ? AND domain = ?", alias, domain).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, dberrors.QueryFailed.Wrap(err, "failed to find paymail by alias and domain")
	}

	return p.newPaymailModel(row), nil
}

// FindForUser returns a paymail by alias and domain for given user.
func (p *Paymails) FindForUser(ctx context.Context, alias, domain, userID string) (*paymailsmodels.Paymail, error) {
	var row database.Paymail
	if err := p.db.
		WithContext(ctx).
		Where("alias = ? AND domain = ? AND user_id = ?", alias, domain, userID).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, dberrors.QueryFailed.Wrap(err, "failed to find paymail by alias, domain and user")
	}

	return p.newPaymailModel(row), nil
}

// GetDefault returns a default paymail for user.
func (p *Paymails) GetDefault(ctx context.Context, userID string) (*paymailsmodels.Paymail, error) {
	var row database.Paymail
	if err := p.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		First(&row).Error; err != nil {
		return nil, dberrors.QueryFailed.Wrap(err, "failed to get default paymail for user")
	}

	return p.newPaymailModel(row), nil
}

func (p *Paymails) newPaymailModel(row database.Paymail) *paymailsmodels.Paymail {
	return &paymailsmodels.Paymail{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,

		Alias:  row.Alias,
		Domain: row.Domain,

		PublicName: row.PublicName,
		Avatar:     row.Avatar,

		UserID: row.UserID,
	}
}
