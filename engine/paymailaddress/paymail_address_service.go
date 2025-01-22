package paymailaddress

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/paerrors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type service struct {
	repo *repository.Paymails
}

// NewService creates a new paymail address service.
func NewService(repo *repository.Paymails) Service {
	return &service{repo: repo}
}

// HasPaymailAddress checks if the given address belongs to a given User.
func (s *service) HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error) {
	alias, domain, _ := paymail.SanitizePaymail(address)
	pm, err := s.repo.Get(ctx, alias, domain)
	if err != nil {
		return false, err
	}

	if pm == nil {
		return false, nil
	}

	return pm.UserID == userID, nil
}

// GetDefaultPaymailAddress returns the default paymail address for the given xPubId.
func (s *service) GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error) {
	pm, err := s.repo.GetDefault(ctx, xPubID)
	if err != nil {
		return "", spverrors.ErrInternal.Wrap(err)
	} else if pm == nil {
		return "", paerrors.ErrNoDefaultPaymailAddress
	}

	return pm.Alias + "@" + pm.Domain, nil
}
