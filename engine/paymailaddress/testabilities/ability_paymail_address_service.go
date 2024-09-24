package testabilities

import (
	"context"
	"slices"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

// PaymailAddressServiceGiven represents the operations that helps with prepare given state of the unit tests environment.
type PaymailAddressServiceGiven interface {
	NewPaymailAddressService() paymailaddress.Service
}

type paymailAddressServiceAbility struct {
	t              testing.TB
	mockRepository *mockRepository
}

// New creates a new test ability.
func New(t testing.TB) (given PaymailAddressServiceGiven) {
	repo := newMockedRepository(t)
	ability := &paymailAddressServiceAbility{
		t:              t,
		mockRepository: repo,
	}

	return ability
}

// NewPaymailAddressService creates a new instance of the paymail address service to use in tests.
func (p *paymailAddressServiceAbility) NewPaymailAddressService() paymailaddress.Service {
	return paymailaddress.NewService(p.mockRepository.getXPubIDByPaymailAddress, p.mockRepository.getPaymailAddressesByXPubIDOrderByCreatedAsc)
}

type mockRepository struct {
	t     testing.TB
	users []fixtures.User
}

func newMockedRepository(t testing.TB) *mockRepository {
	return &mockRepository{
		t:     t,
		users: fixtures.All(),
	}
}

func (s *mockRepository) getXPubIDByPaymailAddress(_ context.Context, paymailAddress string) (string, error) {
	for _, user := range s.users {
		if slices.Contains(user.Paymails, paymailAddress) && user.XPubID != "" {
			return user.XPubID, nil
		}
	}
	return "", spverrors.ErrCouldNotFindPaymail
}

func (s *mockRepository) getPaymailAddressesByXPubIDOrderByCreatedAsc(_ context.Context, xPubID string) ([]string, error) {
	for _, user := range s.users {
		if user.XPubID == xPubID {
			return user.Paymails, nil
		}
	}
	return make([]string, 0), nil
}
