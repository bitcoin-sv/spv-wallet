package users

import (
	"context"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// Service is a user domain service
type Service struct {
	usersRepo UserRepo
	config    *config.AppConfig
}

// NewService creates a new user service
func NewService(users UserRepo, cfg *config.AppConfig) *Service {
	return &Service{
		usersRepo: users,
		config:    cfg,
	}
}

// Create creates a new user
func (s *Service) Create(ctx context.Context, newUser *usersmodels.NewUser) (*usersmodels.User, error) {
	if newUser.Paymail != nil {
		if err := s.config.Paymail.CheckDomain(newUser.Paymail.Domain); err != nil {
			return nil, spverrors.Wrapf(err, "invalid domain during user creation")
		}
		if err := newUser.Paymail.ValidateAvatar(); err != nil {
			return nil, spverrors.Wrapf(err, "invalid avatar url during user creation")
		}
		if newUser.Paymail.PublicName == "" {
			newUser.Paymail.PublicName = newUser.Paymail.Alias
		}
	}
	createdUser, err := s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create user")
	}
	return createdUser, nil
}

// Exists checks if a user exists
func (s *Service) Exists(ctx context.Context, userID string) (bool, error) {
	exists, err := s.usersRepo.Exists(ctx, userID)
	if err != nil {
		return false, spverrors.Wrapf(err, "failed to check if user exists")
	}
	return exists, nil
}

// GetByID returns a user with paymails
func (s *Service) GetByID(ctx context.Context, userID string) (*usersmodels.User, error) {
	user, err := s.usersRepo.Get(ctx, userID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user with paymails")
	}
	return user, nil
}

// GetIDByPubKey returns the user ID selected by pubKey
func (s *Service) GetIDByPubKey(ctx context.Context, pubKey string) (string, error) {
	userID, err := s.usersRepo.GetIDByPubKey(ctx, pubKey)
	if err != nil {
		return "", spverrors.Wrapf(err, "Cannot get user")
	}

	return userID, nil
}

// GetPubKey returns the go-sdk primitives.PublicKey object from the user's PubKey string selected by userID
func (s *Service) GetPubKey(ctx context.Context, userID string) (*primitives.PublicKey, error) {
	user, err := s.usersRepo.Get(ctx, userID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "Cannot get user")
	}

	pubKey, err := user.PubKeyObj()
	if err != nil {
		return nil, spverrors.Wrapf(err, "Cannot get user's public key")
	}
	return pubKey, nil
}

// GetBalance returns current balance for the user
func (s *Service) GetBalance(ctx context.Context, userID string) (bsv.Satoshis, error) {
	balance, err := s.usersRepo.GetBalance(ctx, userID, bucket.BSV)
	if err != nil {
		return 0, spverrors.Wrapf(err, "Cannot get user's balance")
	}
	return balance, nil
}
