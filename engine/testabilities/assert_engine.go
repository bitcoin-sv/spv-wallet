package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/stretchr/testify/require"
)

type EngineAssertions interface {
	ExternalPaymailHost() testpaymail.PaymailExternalAssertions
}

type engineAssertions struct {
	eng         engine.ClientInterface
	paymailMock *paymailmock.PaymailClientMock
	t           testing.TB
	require     *require.Assertions
}

func Then(t testing.TB, engFixture EngineFixture) EngineAssertions {
	fixture := engFixture.(*engineFixture)
	return &engineAssertions{
		eng:         fixture.engine,
		paymailMock: fixture.paymailClient,
		t:           t,
		require:     require.New(t),
	}
}

func (e *engineAssertions) ExternalPaymailHost() testpaymail.PaymailExternalAssertions {
	return testpaymail.Then(e.t, e.paymailMock)
}
