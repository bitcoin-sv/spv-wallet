package user

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"

	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/engine/user/usermodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// Service is a user domain service
type Service struct {
	usersRepo UsersRepo
}

// NewService creates a new user service
func NewService(users UsersRepo) *Service {
	return &Service{
		usersRepo: users,
	}
}

// AppendAddress appends P2PKH address to a user
func (s *Service) AppendAddress(ctx context.Context, userID string, address string, customInstructions bsv.CustomInstructions) error {
	err := s.usersRepo.AppendAddress(ctx, userID, &usermodels.NewAddress{
		Address:            address,
		CustomInstructions: customInstructions,
	})
	if err != nil {
		return spverrors.Wrapf(err, "failed to append address")
	}
	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, newUser *usermodels.NewUser) (*usermodels.User, error) {
	createdUser, err := s.usersRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create user")
	}
	return createdUser, nil
}

// AppendPaymail appends a paymail to a user
func (s *Service) AppendPaymail(ctx context.Context, userID string, newPaymail *usermodels.NewPaymail) (*paymailmodels.Paymail, error) {
	createdPaymail, err := s.usersRepo.AppendPaymail(ctx, userID, newPaymail)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to append paymail")
	}
	return createdPaymail, nil
}

// GetWithPaymails returns a user with paymails
func (s *Service) GetWithPaymails(ctx context.Context, userID string) (*usermodels.User, error) {
	user, err := s.usersRepo.GetWithPaymails(ctx, userID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get user with paymails")
	}
	return user, nil
}
