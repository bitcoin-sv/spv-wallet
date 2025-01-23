package paymails

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
)

// Service for paymails
type Service struct {
	paymailsRepo PaymailRepo
	usersService UsersService
}

// NewService creates a new paymails service
func NewService(paymails PaymailRepo, users UsersService) *Service {
	return &Service{
		paymailsRepo: paymails,
		usersService: users,
	}
}

// Create creates a new paymail attached to a user
func (s *Service) Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error) {
	if exists, err := s.usersService.Exists(ctx, newPaymail.UserID); err != nil {
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

// Get returns a paymail by alias and domain
func (s *Service) Get(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error) {
	paymail, err := s.paymailsRepo.Get(ctx, alias, domain)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get paymail")
	}
	return paymail, nil
}
