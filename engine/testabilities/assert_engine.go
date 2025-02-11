package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type EngineAssertions interface {
	ExternalPaymailHost() PaymailExternalAssertions
}

type UserAssertions interface {
	Balance() BalanceAssertions
}

type BalanceAssertions interface {
	IsEqualTo(expected bsv.Satoshis)
	IsGreaterThanOrEqualTo(expected bsv.Satoshis)
	IsZero()
}

type PaymailExternalAssertions interface {
	ReceivedBeefTransaction(sender, beef, reference string)
}

type PaymailCapabilityCallAssertions interface {
	WithRequestJSONMatching(expectedTemplateFormat string, params map[string]any)
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

func (e *engineAssertions) ExternalPaymailHost() PaymailExternalAssertions {
	return &externalClientAssertions{
		mockPaymail: e.paymailMock,
		t:           e.t,
		require:     e.require,
	}
}
