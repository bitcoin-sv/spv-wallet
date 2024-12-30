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

// GetByPubKey returns a user by its public key. If the user does not exist, it returns error.
func (u *Users) GetByPubKey(ctx context.Context, pubKey string) (*database.User, error) {
	var user database.User
	err := u.db.WithContext(ctx).
		Where("pub_key = ?", pubKey).
		First(&user).Error
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user by public key")
	}

	return &user, nil
}

// GetWithPaymails returns a user by its id with preloaded paymail slist. If the user does not exist, it returns error.
func (u *Users) GetWithPaymails(ctx context.Context, id string) (*database.User, error) {
	var user database.User
	err := u.db.WithContext(ctx).
		Preload("Paymails", func(db *gorm.DB) *gorm.DB {
			//NOTE: To preserve deterministic order necessary to get default paymail as the first one
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user by public key")
	}

	return &user, nil
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

// AppendPaymail appends a paymail to the existing user.
func (u *Users) AppendPaymail(ctx context.Context, userID string, paymailRow *database.Paymail) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.WithContext(ctx).
			Where("id = ?", userID).
			First(&user).Error; err != nil {
			return spverrors.Wrapf(err, "user not found")
		}

		if err := tx.WithContext(ctx).
			Model(&user).
			Association("Paymails").
			Append(paymailRow); err != nil {
			return spverrors.Wrapf(err, "failed to save paymail")
		}

		return nil
	})
	if err != nil {
		return spverrors.Wrapf(err, "failed to append paymail to user")
	}

	return nil
}

// GetBalance returns the balance of a user in a given bucket.
func (u *Users) GetBalance(ctx context.Context, userID string, bucket string) (bsv.Satoshis, error) {
	var balance bsv.Satoshis
	err := u.db.
		WithContext(ctx).
		Model(&database.UserUtxos{}).
		Where("user_id = ? AND bucket = ?", userID, bucket).
		Select("COALESCE(SUM(satoshis), 0)").
		Row().
		Scan(&balance)

	if err != nil {
		return 0, spverrors.Wrapf(err, "failed to get balance")
	}

	return balance, nil
}
