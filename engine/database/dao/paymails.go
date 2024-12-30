package dao

import (
	"errors"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"gorm.io/gorm"
)

// Paymails is a data access object for users.
type Paymails struct {
	db *gorm.DB
}

// NewPaymailsAccessObject creates a new Paymails data access object.
func NewPaymailsAccessObject(db *gorm.DB) *Paymails {
	return &Paymails{db: db}
}

// GetPaymailByAlias returns a paymail by alias and domain.
func (u *Paymails) GetPaymailByAlias(alias, domain string) (*database.Paymail, error) {
	var paymail database.Paymail
	if err := u.db.Preload("User").Where("alias = ? AND domain = ?", alias, domain).First(&paymail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &paymail, nil
}
