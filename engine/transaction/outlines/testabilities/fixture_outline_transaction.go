package testabilities

import (
	"testing"

	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
)

// TransactionOutlineFixture is a test fixture - used for establishing environment for test.
type TransactionOutlineFixture interface {
	NewTransactionOutlinesService() outlines.Service
	ExternalRecipientHost() tpaymail.PaymailHostFixture
}

type transactionOutlineAbility struct {
	t                     testing.TB
	paymailClientAbility  tpaymail.PaymailClientFixture
	paymailAddressService outlines.PaymailAddressService
}

// Given creates a new test fixture.
func Given(t testing.TB) (given TransactionOutlineFixture) {
	ability := &transactionOutlineAbility{
		t:                     t,
		paymailClientAbility:  tpaymail.Given(t),
		paymailAddressService: newPaymailAddressServiceMock(t),
	}
	return ability
}

// ExternalRecipientHost returns test fixture for setting up mocked paymail host.
func (a *transactionOutlineAbility) ExternalRecipientHost() tpaymail.PaymailHostFixture {
	return a.paymailClientAbility.ExternalPaymailHost()
}

// NewTransactionOutlinesService creates a new transaction outline service to use in tests.
func (a *transactionOutlineAbility) NewTransactionOutlinesService() outlines.Service {
	return outlines.NewService(
		a.paymailClientAbility.NewPaymailClientService(),
		a.paymailAddressService,
		tester.Logger(a.t),
	)
}
