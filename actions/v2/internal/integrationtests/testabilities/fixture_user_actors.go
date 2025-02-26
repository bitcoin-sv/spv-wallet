package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

type IntegrationTestFixtures interface {
	StartedSPVWalletV2(opts ...testengine.ConfigOpts) (cleanup func())
	Paymail() testpaymail.PaymailClientFixture
	Alice() *fixtures.User
	Bob() *fixtures.User
}

type fixture struct {
	testabilities.SPVWalletApplicationFixture
	t       testing.TB
	alice   *user
	bob     *user
	charlie *user
}

func newFixture(t testing.TB, appFixture testabilities.SPVWalletApplicationFixture) *fixture {
	return &fixture{
		t:                           t,
		SPVWalletApplicationFixture: appFixture,
		alice: &user{
			User: fixtures.Sender,
			app:  appFixture,
			t:    t,
		},
		bob: &user{
			User: fixtures.RecipientInternal,
			app:  appFixture,
			t:    t,
		},
		charlie: &user{
			User: fixtures.RecipientExternal,
			app:  appFixture,
			t:    t,
		},
	}
}

func (f *fixture) StartedSPVWalletV2(opts ...testengine.ConfigOpts) (cleanup func()) {
	cleanup = f.StartedSPVWalletWithConfiguration(append(opts, testengine.WithV2())...)
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
