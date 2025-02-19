package testabilities

import (
	"context"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// TransactionOutlineFixture is a test fixture - used for establishing environment for test.
type TransactionOutlineFixture interface {
	MinimumValidTransactionSpec() *outlines.TransactionSpec
	NewTransactionOutlinesService() outlines.Service
	ExternalRecipientHost() tpaymail.PaymailHostFixture
	UserHasNotEnoughFunds()
	UTXOSelector() UTXOSelectorFixture
}

type transactionOutlineAbility struct {
	t                     testing.TB
	paymailClientAbility  tpaymail.PaymailClientFixture
	paymailAddressService outlines.PaymailAddressService
	utxoSelector          mockedUTXOSelector
	feeUnit               bsv.FeeUnit
}

func (a *transactionOutlineAbility) MinimumValidTransactionSpec() *outlines.TransactionSpec {
	return &outlines.TransactionSpec{
		UserID: fixtures.Sender.ID(),
		Outputs: outlines.NewOutputsSpecs(&outlines.OpReturn{
			Data: []string{"test"},
		}),
	}
}

// Given creates a new test fixture.
func Given(t testing.TB) (given TransactionOutlineFixture) {
	ability := &transactionOutlineAbility{
		t:                     t,
		paymailClientAbility:  tpaymail.Given(t),
		paymailAddressService: newPaymailAddressServiceMock(t),
		feeUnit:               bsv.FeeUnit{Satoshis: 1, Bytes: 1000},
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
		&a.utxoSelector,
		a.feeUnit,
		tester.Logger(a.t),
		pubKeyGetter{},
	)
}

func (a *transactionOutlineAbility) UTXOSelector() UTXOSelectorFixture {
	return &a.utxoSelector
}

func (a *transactionOutlineAbility) UserHasNotEnoughFunds() {
	a.utxoSelector.WillReturnNoUTXOs()
}

type pubKeyGetter struct{}

func (p pubKeyGetter) GetPubKey(ctx context.Context, _ string) (*ec.PublicKey, error) {
	return fixtures.Sender.PublicKey(), nil
}
