package testabilities

import (
	"testing"

	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	tpaymailaddress "github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
)

// DraftTransactionFixture is a test fixture - used for establishing environment for test.
type DraftTransactionFixture interface {
	NewDraftTransactionService() draft.Service
	ExternalRecipientHost() tpaymail.PaymailHostFixture
}

type draftTransactionAbility struct {
	t                     testing.TB
	paymailClientAbility  tpaymail.PaymailClientFixture
	paymailAddressAbility tpaymailaddress.PaymailAddressServiceFixture
}

// Given creates a new test fixture.
func Given(t testing.TB) (given DraftTransactionFixture) {
	ability := &draftTransactionAbility{
		t:                     t,
		paymailClientAbility:  tpaymail.Given(t),
		paymailAddressAbility: tpaymailaddress.Given(t),
	}
	return ability
}

// ExternalRecipientHost returns test fixture for setting up mocked paymail host.
func (a *draftTransactionAbility) ExternalRecipientHost() tpaymail.PaymailHostFixture {
	return a.paymailClientAbility.ExternalPaymailHost()
}

// NewDraftTransactionService creates a new draft transaction service to use in tests.
func (a *draftTransactionAbility) NewDraftTransactionService() draft.Service {
	return draft.NewDraftService(a.paymailClientAbility.NewPaymailClientService(), a.paymailAddressAbility.NewPaymailAddressService(), tester.Logger(a.t))
}
