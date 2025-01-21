package repository

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
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
func (p *Paymails) Create(ctx context.Context, newPaymail *domainmodels.NewPaymail) (*domainmodels.Paymail, error) {
	row := database.Paymail{
		Alias:  newPaymail.Alias,
		Domain: newPaymail.Domain,

		PublicName: newPaymail.PublicName,
		Avatar:     newPaymail.Avatar,

		UserID: newPaymail.UserID,
	}

	if err := p.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}

	return &domainmodels.Paymail{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,

		Alias:  row.Alias,
		Domain: row.Domain,

		PublicName: row.PublicName,
		Avatar:     row.Avatar,

		UserID: row.UserID,
	}, nil
}

// Get returns a paymail by alias and domain.
func (p *Paymails) Get(ctx context.Context, alias, domain string) (*domainmodels.Paymail, error) {
	var row database.Paymail
	if err := p.db.
		WithContext(ctx).
		Where("alias = ? AND domain = ?", alias, domain).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domainmodels.Paymail{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,

		Alias:  row.Alias,
		Domain: row.Domain,

		PublicName: row.PublicName,
		Avatar:     row.Avatar,

		UserID: row.UserID,
	}, nil
}
