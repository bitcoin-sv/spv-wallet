package testabilities

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type EngineAssertions interface {
	User(fixtures.User) UserAssertions
	PaymailClient() PaymailClientAssertions
}

type UserAssertions interface {
	Balance() BalanceAssertions
}

type BalanceAssertions interface {
	IsEqualTo(expected bsv.Satoshis)
	IsGreaterThanOrEqualTo(expected bsv.Satoshis)
	IsZero()
}

type PaymailClientAssertions interface {
	ExternalPaymailHost() PaymailExternalAssertions
}

type PaymailExternalAssertions interface {
	Called(urlRegex string) PaymailCapabilityCallAssertions
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

type userAssertions struct {
	eng     engine.ClientInterface
	t       testing.TB
	user    fixtures.User
	require *require.Assertions
}

func (e *engineAssertions) User(user fixtures.User) UserAssertions {
	return &userAssertions{
		eng:     e.eng,
		t:       e.t,
		user:    user,
		require: e.require,
	}
}

func (u *userAssertions) Balance() BalanceAssertions {
	return u
}

func (u *userAssertions) balance() bsv.Satoshis {
	u.t.Helper()
	actual, err := u.eng.UsersService().GetBalance(context.Background(), u.user.ID())
	u.require.NoError(err)
	return actual
}

func (u *userAssertions) IsEqualTo(expected bsv.Satoshis) {
	u.t.Helper()
	actual := u.balance()
	require.Equal(u.t, expected, actual)
}

func (u *userAssertions) IsGreaterThanOrEqualTo(expected bsv.Satoshis) {
	u.t.Helper()
	actual := u.balance()
	require.GreaterOrEqual(u.t, actual, expected)
}

func (u *userAssertions) IsZero() {
	u.t.Helper()
	actual := u.balance()
	require.Zero(u.t, actual)
}

func (e *engineAssertions) PaymailClient() PaymailClientAssertions {
	return e
}

func (e *engineAssertions) ExternalPaymailHost() PaymailExternalAssertions {
	return &externalClientAssertions{
		mockPaymail: e.paymailMock,
		t:           e.t,
		require:     e.require,
	}
}

type externalClientAssertions struct {
	t           testing.TB
	require     *require.Assertions
	mockPaymail *paymailmock.PaymailClientMock
}

func (e *externalClientAssertions) Called(urlRegex string) PaymailCapabilityCallAssertions {
	details := e.mockPaymail.GetCallByRegex(urlRegex)
	e.require.NotNil(details, "Expected call to %s", urlRegex)
	return &paymailCapabilityCallAssertions{
		t:           e.t,
		callDetails: details,
	}
}

type paymailCapabilityCallAssertions struct {
	t           testing.TB
	callDetails *paymailmock.CallDetails
}

func (p *paymailCapabilityCallAssertions) WithRequestJSONMatching(expectedTemplateFormat string, params map[string]any) {
	jsonrequire.Match(p.t, expectedTemplateFormat, params, string(p.callDetails.RequestBody))
}
