package repository

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/user/usermodels"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/datatypes"
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
func (u *Users) GetWithPaymails(ctx context.Context, userID string) (*usermodels.User, error) {
	var user database.User
	err := u.db.WithContext(ctx).
		Preload("Paymails", func(db *gorm.DB) *gorm.DB {
			//NOTE: To preserve deterministic order necessary to get default paymail as the first one
			return db.Order("created_at ASC")
		}).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user by public key")
	}

	return &usermodels.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		PublicKey: user.PubKey,
		Paymails: utils.MapSlice(user.Paymails, func(p *database.Paymail) *paymailmodels.Paymail {
			return &paymailmodels.Paymail{
				ID:        p.ID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,

				Alias:  p.Alias,
				Domain: p.Domain,

				PublicName: p.PublicName,
				Avatar:     p.Avatar,

				UserID: p.UserID,
				User:   p.User,
			}
		}),
	}, nil
}

// CreateUser saves new user to the database.
func (u *Users) CreateUser(ctx context.Context, newUser *usermodels.NewUser) (*usermodels.User, error) {
	query := u.db.WithContext(ctx)

	row := &database.User{
		PubKey: newUser.PublicKey,
	}
	if newUser.Paymail != nil {
		row.Paymails = []*database.Paymail{{
			Alias:      newUser.Paymail.Alias,
			Domain:     newUser.Paymail.Domain,
			PublicName: newUser.Paymail.PublicName,
			Avatar:     newUser.Paymail.Avatar,
		}}
	}

	if err := query.Create(row).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to save user")
	}

	return &usermodels.User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		PublicKey: row.PubKey,
		Paymails: utils.MapSlice(row.Paymails, func(p *database.Paymail) *paymailmodels.Paymail {
			return &paymailmodels.Paymail{
				ID:        p.ID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,

				Alias:  p.Alias,
				Domain: p.Domain,

				PublicName: p.PublicName,
				Avatar:     p.Avatar,

				UserID: p.UserID,
				User:   p.User,
			}
		}),
	}, nil
}

// AppendAddress appends an address to the database.
func (u *Users) AppendAddress(ctx context.Context, userID string, newAddress *usermodels.NewAddress) error {
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.
			Where("id = ?", userID).
			First(&user).Error; err != nil {
			return spverrors.Wrapf(err, "user not found")
		}

		if err := tx.
			Model(&user).
			Association("Addresses").
			Append(&database.Address{
				Address:            newAddress.Address,
				CustomInstructions: datatypes.NewJSONSlice(newAddress.CustomInstructions),
			}); err != nil {
			return spverrors.Wrapf(err, "failed to save address")
		}

		return nil
	})
	if err != nil {
		return spverrors.Wrapf(err, "failed to append address to user")
	}

	return nil
}

// AppendPaymail appends a paymail to the existing user.
func (u *Users) AppendPaymail(ctx context.Context, userID string, newPaymail *usermodels.NewPaymail) (*paymailmodels.Paymail, error) {
	paymailRow := &database.Paymail{
		Alias:      newPaymail.Alias,
		Domain:     newPaymail.Domain,
		PublicName: newPaymail.PublicName,
		Avatar:     newPaymail.Avatar,
	}
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.
			Where("id = ?", userID).
			First(&user).Error; err != nil {
			return spverrors.Wrapf(err, "user not found")
		}

		if err := tx.
			Model(&user).
			Association("Paymails").
			Append(paymailRow); err != nil {
			return spverrors.Wrapf(err, "failed to save paymail")
		}

		return nil
	})
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to append paymail to user")
	}

	return &paymailmodels.Paymail{
		ID:        paymailRow.ID,
		CreatedAt: paymailRow.CreatedAt,
		UpdatedAt: paymailRow.UpdatedAt,

		Alias:  paymailRow.Alias,
		Domain: paymailRow.Domain,

		PublicName: paymailRow.PublicName,
		Avatar:     paymailRow.Avatar,

		UserID: paymailRow.UserID,
		User:   paymailRow.User,
	}, nil
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
