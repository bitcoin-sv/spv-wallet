package repository

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	dberrors "github.com/bitcoin-sv/spv-wallet/engine/v2/database/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/samber/lo"
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

// Exists checks if a user exists in the database.
func (u *Users) Exists(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&database.User{}).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, dberrors.QueryFailed.Wrap(err, "failed to check if user exists by ID")
	}

	return count > 0, nil
}

// GetIDByPubKey returns a user by its public key. If the user does not exist, it returns error.
func (u *Users) GetIDByPubKey(ctx context.Context, pubKey string) (string, error) {
	var user struct {
		ID string
	}
	err := u.db.WithContext(ctx).
		Model(&database.User{}).
		Where("pub_key = ?", pubKey).
		First(&user).Error
	if err != nil {
		return "", dberrors.QueryFailed.Wrap(err, "failed to get user by public key")
	}

	return user.ID, nil
}

// Get returns a user by its id with preloaded paymail slist. If the user does not exist, it returns error.
func (u *Users) Get(ctx context.Context, userID string) (*usersmodels.User, error) {
	var user database.User
	err := u.db.WithContext(ctx).
		Scopes(withPaymailsScope).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, dberrors.QueryOrNotFoundError(err, "failed to get user by ID")
	}

	return mapToDomainUser(&user), nil
}

// Create saves new user to the database.
func (u *Users) Create(ctx context.Context, newUser *usersmodels.NewUser) (*usersmodels.User, error) {
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
		return nil, dberrors.QueryFailed.Wrap(err, "failed to create new user")
	}

	return mapToDomainUser(row), nil
}

// GetBalance returns the balance of a user in a given bucket.
func (u *Users) GetBalance(ctx context.Context, userID string, bucket bucket.Name) (bsv.Satoshis, error) {
	var balance bsv.Satoshis
	err := u.db.
		WithContext(ctx).
		Model(&database.UserUTXO{}).
		Where("user_id = ? AND bucket = ?", userID, bucket).
		Select("COALESCE(SUM(satoshis), 0)").
		Row().
		Scan(&balance)

	if err != nil {
		return 0, dberrors.QueryFailed.Wrap(err, "failed to get balance for user by ID")
	}

	return balance, nil
}

func mapToDomainUser(user *database.User) *usersmodels.User {
	return &usersmodels.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		PublicKey: user.PubKey,
		Paymails: lo.Map(user.Paymails, func(p *database.Paymail, _ int) *paymailsmodels.Paymail {
			return &paymailsmodels.Paymail{
				ID:        p.ID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,

				Alias:  p.Alias,
				Domain: p.Domain,

				PublicName: p.PublicName,
				Avatar:     p.Avatar,

				UserID: p.UserID,
			}
		}),
	}
}

func withPaymailsScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Paymails", func(db *gorm.DB) *gorm.DB {
		// NOTE: To preserve deterministic order necessary to get default paymail as the first one
		return db.Order("created_at ASC")
	})
}
