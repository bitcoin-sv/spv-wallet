package users

import (
	"context"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Service is a user domain service
type Service struct {
	usersRepo     UserRepo
	addressesRepo AddressRepo
	paymailsRepo  PaymailRepo
}

// UserService creates a new user service
func UserService(users UserRepo, addresses AddressRepo, paymailsRepo PaymailRepo) *Service {
	return &Service{
		usersRepo:     users,
		addressesRepo: addresses,
		paymailsRepo:  paymailsRepo,
	}
}

// AppendAddress appends P2PKH address to a user
func (s *Service) AppendAddress(ctx context.Context, newAddress domainmodels.NewAddress) error {
	if exists, err := s.usersRepo.Exists(ctx, newAddress.UserID); err != nil {
		return spverrors.Wrapf(err, "failed to check if user exists")
	} else if !exists {
		return spverrors.Newf("user does not exist")
	}

	err := s.addressesRepo.Create(ctx, &newAddress)
	if err != nil {
		return spverrors.Wrapf(err, "failed to append address")
	}
	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, newUser *domainmodels.NewUser) (*domainmodels.User, error) {
	createdUser, err := s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create user")
	}
	return createdUser, nil
}

// AppendPaymail appends a paymail to a user
func (s *Service) AppendPaymail(ctx context.Context, newPaymail *domainmodels.NewPaymail) (*domainmodels.Paymail, error) {
	if exists, err := s.usersRepo.Exists(ctx, newPaymail.UserID); err != nil {
		return nil, spverrors.Wrapf(err, "failed to check if user exists")
	} else if !exists {
		return nil, spverrors.Newf("user does not exist")
	}

	createdPaymail, err := s.paymailsRepo.Create(ctx, newPaymail)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to append paymail")
	}
	return createdPaymail, nil
}

// GetByID returns a user with paymails
func (s *Service) GetByID(ctx context.Context, userID string) (*domainmodels.User, error) {
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
