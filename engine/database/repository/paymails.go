package repository

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
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

// Find returns a paymail by alias and domain.
func (p *Paymails) Find(ctx context.Context, alias, domain string) (*database.Paymail, error) {
	var paymail database.Paymail
	if err := p.db.
		WithContext(ctx).
		Preload("User").
		Where("alias = ? AND domain = ?", alias, domain).
		First(&paymail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &paymail, nil
}

// FindForUser returns a paymail by alias and domain for given user.
func (p *Paymails) FindForUser(ctx context.Context, alias, domain, userID string) (*database.Paymail, error) {
	var paymail database.Paymail
	if err := p.db.
		WithContext(ctx).
		Preload("User").
		Where("alias = ? AND domain = ? AND user_id = ?", alias, domain, userID).
		First(&paymail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &paymail, nil
}

// GetDefault returns a default paymail for user.
func (p *Paymails) GetDefault(ctx context.Context, userID string) (*database.Paymail, error) {
	var paymail database.Paymail
	if err := p.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		First(&paymail).Error; err != nil {
		return nil, err
	}

	return &paymail, nil
}
