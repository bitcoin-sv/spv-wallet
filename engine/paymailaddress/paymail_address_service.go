package paymailaddress

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/paerrors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
)

type service struct {
	repo PaymailRepo
}

// NewService creates a new paymail address service.
func NewService(repo PaymailRepo) Service {
	return &service{repo: repo}
}

// HasPaymailAddress checks if the given address belongs to a given User.
func (s *service) HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error) {
	alias, domain, sanitized := paymail.SanitizePaymail(address)
	if sanitized == "" {
		return false, paerrors.ErrInvalidPaymailAddress
	}
	pm, err := s.repo.FindForUser(ctx, alias, domain, userID)
	if err != nil {
		return false, spverrors.ErrInternal.Wrap(err)
	}

	return pm != nil, nil
}

// GetDefaultPaymailAddress returns the default paymail address for the given xPubId.
func (s *service) GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error) {
	pm, err := s.repo.GetDefault(ctx, xPubID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", paerrors.ErrNoDefaultPaymailAddress
	} else if err != nil {
		return "", spverrors.ErrInternal.Wrap(err)
	}

	return pm.Alias + "@" + pm.Domain, nil
}
