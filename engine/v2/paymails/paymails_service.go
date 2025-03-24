package paymails

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/joomcode/errorx"
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
		return nil, errorx.Decorate(err, "invalid domain during paymail creation")
	}

	if exists, err := s.usersService.Exists(ctx, newPaymail.UserID); err != nil {
		return nil, errorx.Decorate(err, "paymails service failed to check if user exists")
	} else if !exists {
		return nil, paymailerrors.UserDoesntExist.New("user with ID %s does not exist", newPaymail.UserID)
	}

	if err := newPaymail.ValidateAvatar(); err != nil {
		return nil, errorx.Decorate(err, "invalid avatar url during user creation")
	}
	if newPaymail.PublicName == "" {
		newPaymail.PublicName = newPaymail.Alias
	}

	createdPaymail, err := s.paymailsRepo.Create(ctx, newPaymail)
	if err != nil {
		return nil, errorx.Decorate(err, "paymails service failed to create new paymail for user")
	}
	return createdPaymail, nil
}

// Find returns a paymail by alias and domain
func (s *Service) Find(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error) {
	address, err := s.paymailsRepo.Find(ctx, alias, domain)
	if err != nil {
		return nil, errorx.Decorate(err, "paymails service failed to find paymail by alias and domain")
	}
	return address, nil
}

// HasPaymailAddress checks if the given address belongs to a given User.
func (s *Service) HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error) {
	alias, domain, sanitized := paymail.SanitizePaymail(address)
	if sanitized == "" {
		return false, paymailerrors.InvalidPaymailAddress.New("invalid paymail address: %s", address)
	}
	pm, err := s.paymailsRepo.FindForUser(ctx, alias, domain, userID)
	if err != nil {
		return false, errorx.Decorate(err, "paymails service failed to find paymail by alias and domain for user")
	}

	return pm != nil, nil
}

// GetDefaultPaymailAddress returns the default paymail address for the given xPubId.
func (s *Service) GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error) {
	pm, err := s.paymailsRepo.GetDefault(ctx, xPubID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", paymailerrors.NoDefaultPaymailAddress.New("no default paymail address found for xPubID: %s", xPubID)
	} else if err != nil {
		return "", errorx.Decorate(err, "paymails service failed to get default paymail address for xPubID: %s", xPubID)
	}

	return pm.Alias + "@" + pm.Domain, nil
}
