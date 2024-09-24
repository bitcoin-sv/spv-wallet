package testabilities

import (
	"testing"

	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	tpaymailaddress "github.com/bitcoin-sv/spv-wallet/engine/paymailaddress/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// DraftTransactionGiven represents the operations that helps with prepare given state of the unit tests environment.
type DraftTransactionGiven interface {
	NewDraftTransactionService() draft.Service
	ExternalRecipientHost() tpaymail.PaymailHostGiven
}

// DraftTransactionThen represents the assertions on the result of the operation.
type DraftTransactionThen interface {
	Created(transaction *draft.Transaction) DraftTransactionAssertion
}

type DraftTransactionAssertion interface {
	WithNoError(err error) DraftTransactionAssertion
	HasOutputs(count int) DraftTransactionAssertion
	Output(index int) DraftTransactionOutputAssertion
}

type DraftTransactionOutputAssertion interface {
	HasBucket(bucket transaction.Bucket) DraftTransactionOutputAssertion
	HasSatoshis(satoshis bsv.Satoshis) DraftTransactionOutputAssertion
	HasLockingScript(lockingScript string) DraftTransactionOutputAssertion
	IsPaymail() DraftTransactionPaymailOutputAssertion
}

type DraftTransactionPaymailOutputAssertion interface {
	HasReceiver(receiver string) DraftTransactionPaymailOutputAssertion
	HasSender(sender string) DraftTransactionPaymailOutputAssertion
	HasReference(reference string) DraftTransactionPaymailOutputAssertion
	And() DraftTransactionAssertion
}

type draftTransactionAbility struct {
	t                     testing.TB
	paymailClientAbility  tpaymail.PaymailClientGiven
	paymailAddressAbility tpaymailaddress.PaymailAddressServiceGiven
}

// New creates a test ability.
func New(t testing.TB) (given DraftTransactionGiven, then DraftTransactionThen) {
	ability := &draftTransactionAbility{
		t:                     t,
		paymailClientAbility:  tpaymail.New(t),
		paymailAddressAbility: tpaymailaddress.New(t),
	}
	return ability, nil
}

// Then creates a new then part of test ability.
func Then(t testing.TB) (then DraftTransactionThen) {
	return nil
}

// Given creates a new given part of test ability.
func Given(t testing.TB) (given DraftTransactionGiven) {
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
