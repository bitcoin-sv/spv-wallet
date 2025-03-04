package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
)

type IntegrationTestFixtures interface {
	StartedSPVWalletV2() (cleanup func())
	Paymail() testpaymail.PaymailClientFixture

	Alice() *fixtures.User
	Bob() *fixtures.User
	Charlie() *fixtures.User
}

type fixture struct {
	testabilities.SPVWalletApplicationFixture
	t       testing.TB
	alice   *user
	bob     *user
	charlie *user
}

func newFixture(t testing.TB, appFixture testabilities.SPVWalletApplicationFixture) *fixture {
	txFixture := txtestability.Given(t)

	return &fixture{
		t:                           t,
		SPVWalletApplicationFixture: appFixture,
		alice: &user{
			User:      fixtures.Sender,
			app:       appFixture,
			txFixture: txFixture,
			t:         t,
		},
		bob: &user{
			User:      fixtures.RecipientInternal,
			app:       appFixture,
			txFixture: txFixture,
			t:         t,
		},
		charlie: &user{
			User:      fixtures.RecipientExternal,
			app:       appFixture,
			txFixture: txFixture,
			t:         t,
		},
	}
}

func (f *fixture) StartedSPVWalletV2() (cleanup func()) {
	cleanup = f.StartedSPVWalletWithConfiguration(testengine.WithV2())
	f.Paymail().ExternalPaymailHost().WillRespondWithP2PWithBEEFCapabilities()
	return
}

func (f *fixture) Alice() *fixtures.User {
	alice := f.alice.User
	return &alice
}

func (f *fixture) Bob() *fixtures.User {
	bob := f.bob.User
	return &bob
}

func (f *fixture) Charlie() *fixtures.User {
	charlie := f.charlie.User
	return &charlie
}
