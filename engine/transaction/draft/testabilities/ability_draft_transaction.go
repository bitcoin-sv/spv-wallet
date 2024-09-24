package testabilities

import (
	"testing"

	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	tpaymailaddress "github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
)

// DraftTransactionGiven represents the operations that helps with prepare given state of the unit tests environment.
type DraftTransactionGiven interface {
	NewDraftTransactionService() draft.Service
	ExternalRecipientHost() tpaymail.PaymailHostGiven
}

type draftTransactionAbility struct {
	t                     testing.TB
	paymailClientAbility  tpaymail.PaymailClientGiven
	paymailAddressAbility tpaymailaddress.PaymailAddressServiceGiven
}

// New creates a new test ability.
func New(t testing.TB) (given DraftTransactionGiven) {
	ability := &draftTransactionAbility{
		t:                     t,
		paymailClientAbility:  tpaymail.New(t),
		paymailAddressAbility: tpaymailaddress.New(t),
	}
	return ability
}

// ExternalRecipientHost returns helper for setting up mocked paymail host.
func (a *draftTransactionAbility) ExternalRecipientHost() tpaymail.PaymailHostGiven {
	return a.paymailClientAbility.ExternalPaymailHost()
}

// NewDraftTransactionService creates a new draft transaction service to use in tests.
func (a *draftTransactionAbility) NewDraftTransactionService() draft.Service {
	return draft.NewDraftService(a.paymailClientAbility.NewPaymailClientService(), a.paymailAddressAbility.NewPaymailAddressService(), tester.Logger(a.t))
}
