package paymails

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"gorm.io/gorm"
)

// Service for paymails
type Service struct {
	paymailsRepo PaymailRepo
	usersService UsersService
	config       *config.AppConfig
}

// NewService creates a new paymails service
func NewService(paymails PaymailRepo, users UsersService, cfg *config.AppConfig) *Service {
	return &Service{
		paymailsRepo: paymails,
		usersService: users,
		config:       cfg,
	}
}

// Create creates a new paymail attached to a user
func (s *Service) Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error) {
	if err := s.config.Paymail.CheckDomain(newPaymail.Domain); err != nil {
		return nil, spverrors.Wrapf(err, "invalid domain during paymail creation")
	}

	if exists, err := s.usersService.Exists(ctx, newPaymail.UserID); err != nil {
		return nil, spverrors.Wrapf(err, "failed to check if user exists")
	} else if !exists {
		return nil, spverrors.Newf("user does not exist")
	}

	if err := newPaymail.ValidateAvatar(); err != nil {
		return nil, spverrors.Newf("invalid avatar url during paymail creation")
	}
	if newPaymail.PublicName == "" {
		newPaymail.PublicName = newPaymail.Alias
	}

	createdPaymail, err := s.paymailsRepo.Create(ctx, newPaymail)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to append paymail")
	}
	return createdPaymail, nil
}

// Find returns a paymail by alias and domain
func (s *Service) Find(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error) {
	paymail, err := s.paymailsRepo.Find(ctx, alias, domain)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get paymail")
	}
	return paymail, nil
}

// HasPaymailAddress checks if the given address belongs to a given User.
func (s *Service) HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error) {
	alias, domain, sanitized := paymail.SanitizePaymail(address)
	if sanitized == "" {
		return false, paymailerrors.ErrInvalidPaymailAddress
	}
	pm, err := s.paymailsRepo.FindForUser(ctx, alias, domain, userID)
	if err != nil {
		return false, spverrors.ErrInternal.Wrap(err)
	}

	return pm != nil, nil
}

// GetDefaultPaymailAddress returns the default paymail address for the given xPubId.
func (s *Service) GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error) {
	pm, err := s.paymailsRepo.GetDefault(ctx, xPubID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", paymailerrors.ErrNoDefaultPaymailAddress
	} else if err != nil {
		return "", spverrors.ErrInternal.Wrap(err)
	}

	return pm.Alias + "@" + pm.Domain, nil
}
