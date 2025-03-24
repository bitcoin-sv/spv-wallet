package testabilities

import (
	"context"
	"slices"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
)

type mockPaymailAddressService struct {
	t     testing.TB
	users []fixtures.User
}

func newPaymailAddressServiceMock(t testing.TB) *mockPaymailAddressService {
	return &mockPaymailAddressService{
		t:     t,
		users: fixtures.InternalUsers(),
	}
}

func (m *mockPaymailAddressService) HasPaymailAddress(_ context.Context, userID string, address string) (bool, error) {
	for _, user := range m.users {
		if user.ID() == userID {
			return slices.Contains(user.Paymails, fixtures.Paymail(address)), nil
		}
	}
	return false, nil
}

func (m *mockPaymailAddressService) GetDefaultPaymailAddress(_ context.Context, userID string) (string, error) {
	for _, user := range m.users {
		if user.ID() == userID && user.DefaultPaymail().Address() != "" {
			return user.DefaultPaymail().Address(), nil
		}
	}
	return "", paymailerrors.NoDefaultPaymailAddress.NewWithNoMessage()
}
