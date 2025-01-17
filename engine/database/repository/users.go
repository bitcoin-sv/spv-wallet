package repository

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// Users is a repository for users.
type Users struct {
	db *gorm.DB
}

// NewUsersRepo creates a new repository for users.
func NewUsersRepo(db *gorm.DB) *Users {
	return &Users{db: db}
}

// Save saves a user to the database.
func (u *Users) Save(ctx context.Context, userRow *database.User) error {
	query := u.db.WithContext(ctx)

	if err := query.Create(userRow).Error; err != nil {
		return spverrors.Wrapf(err, "failed to save user")
	}

	return nil
}

// AppendAddress appends an address to the database.
func (u *Users) AppendAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error {
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

// GetBalance returns the balance of a user in a given bucket.
func (u *Users) GetBalance(ctx context.Context, userID string, bucket string) (bsv.Satoshis, error) {
	var balance bsv.Satoshis
	err := u.db.
		WithContext(ctx).
		Model(&database.UsersUTXO{}).
		Where("user_id = ? AND bucket = ?", userID, bucket).
		Select("COALESCE(SUM(satoshis), 0)").
		Row().
		Scan(&balance)

	if err != nil {
		return 0, spverrors.Wrapf(err, "failed to get balance")
	}

	return balance, nil
}
