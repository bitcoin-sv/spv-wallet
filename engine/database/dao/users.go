package dao

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Users is a data access object for users.
type Users struct {
	db *gorm.DB
}

// NewUsersAccessObject creates a new access object for users.
func NewUsersAccessObject(db *gorm.DB) *Users {
	return &Users{db: db}
}

// SaveUser saves a user to the database.
func (u *Users) SaveUser(ctx context.Context, userRow *database.User) error {
	query := u.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		})

	if err := query.Create(userRow).Error; err != nil {
		return spverrors.Wrapf(err, "failed to save user")
	}

	return nil
}

// GetPaymailByAlias returns a paymail by alias and domain.
func (u *Users) GetPaymailByAlias(alias, domain string) (*database.Paymail, error) {
	var paymail database.Paymail
	if err := u.db.Preload("User").Where("alias = ? AND domain = ?", alias, domain).First(&paymail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &paymail, nil
}

func (u *Users) SaveAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error {
	err := u.db.
		WithContext(ctx).
		Model(userRow).
		Association("Addresses").
		Append(addressRow)

	if err != nil {
		return spverrors.Wrapf(err, "failed to save address")
	}

	return nil
}
