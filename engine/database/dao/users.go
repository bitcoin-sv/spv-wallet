package dao

import (
	"context"

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
