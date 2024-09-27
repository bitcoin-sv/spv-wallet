package paymailaddress

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/paerrors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type service struct {
	// INFO: this is the first step to slowly move paymail address to a separate package.
	getXPubIDByPaymailAddress                    func(ctx context.Context, paymailAddress string) (string, error)
	getPaymailAddressesByXPubIDOrderByCreatedAsc func(ctx context.Context, xPubId string) ([]string, error)
}

// NewService creates a new paymail address service.
func NewService(
	getXPubIDByPaymailAddress func(ctx context.Context, paymailAddress string) (string, error),
	getPaymailAddressesByXPubIDOrderByCreatedAsc func(ctx context.Context, xPubId string) ([]string, error),
) Service {
	if getXPubIDByPaymailAddress == nil {
		panic("getXPubIDByPaymailAddress is required to create paymail address service")
	}
	if getPaymailAddressesByXPubIDOrderByCreatedAsc == nil {
		panic("getPaymailAddressesByXPubIDOrderByCreatedAsc is required to create paymail address service")
	}

	return &service{
		getXPubIDByPaymailAddress:                    getXPubIDByPaymailAddress,
		getPaymailAddressesByXPubIDOrderByCreatedAsc: getPaymailAddressesByXPubIDOrderByCreatedAsc,
	}
}

// HasPaymailAddress checks if the given address belongs to a given xPubId.
func (s *service) HasPaymailAddress(ctx context.Context, xPubID string, address string) (bool, error) {
	paymailXpubID, err := s.getXPubIDByPaymailAddress(ctx, address)
	if err != nil {
		return false, err
	}
	return paymailXpubID == xPubID, nil
}

// GetDefaultPaymailAddress returns the default paymail address for the given xPubId.
func (s *service) GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error) {
	addresses, err := s.getPaymailAddressesByXPubIDOrderByCreatedAsc(ctx, xPubID)
	if err != nil {
		return "", spverrors.ErrInternal.Wrap(err)
	}
	if len(addresses) == 0 {
		return "", paerrors.ErrNoDefaultPaymailAddress
	}
	return addresses[0], nil
}
