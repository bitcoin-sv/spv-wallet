package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/stretchr/testify/require"
)

type EngineAssertions interface {
	ExternalPaymailHost() testpaymail.PaymailExternalAssertions
	ARC() ARCAssertions
}

type ARCAssertions interface {
	Broadcasted() ARCBroadcastAssertions
}

type ARCBroadcastAssertions interface {
	WithTxID(txID string) ARCBroadcastAssertions
	WithCallbackURL(url string) ARCBroadcastAssertions
	WithCallbackToken(token string) ARCBroadcastAssertions
}

type engineAssertions struct {
	eng         engine.ClientInterface
	paymailMock *paymailmock.PaymailClientMock
	t           testing.TB
	require     *require.Assertions
	sniffer     *tester.HTTPSniffer
}

func Then(t testing.TB, engFixture EngineFixture) EngineAssertions {
	fixture := engFixture.(*engineFixture)
	return &engineAssertions{
		eng:         fixture.engine,
		paymailMock: fixture.paymailClient,
		t:           t,
		require:     require.New(t),
		sniffer:     fixture.externalTransportWithSniffer,
	}
}

func (e *engineAssertions) ExternalPaymailHost() testpaymail.PaymailExternalAssertions {
	return testpaymail.Then(e.t, e.paymailMock)
}

func (e *engineAssertions) ARC() ARCAssertions {
	return e
}

func (e *engineAssertions) Broadcasted() ARCBroadcastAssertions {
	e.t.Helper()

	details := e.sniffer.GetCallByRegex("arc.*\\/v1\\/tx")
	e.require.NotNil(details, "Expected call to arc /v1/tx")
	e.require.Equal("POST", details.RequestMethod)
	e.require.Equal(200, details.ResponseCode)

	return &arcBroadcastAssertions{
		details: details,
		t:       e.t,
		require: e.require,
	}
}

type arcBroadcastAssertions struct {
	details *tester.CallDetails
	t       testing.TB
	require *require.Assertions
}

func (a *arcBroadcastAssertions) WithTxID(txID string) ARCBroadcastAssertions {
	rawTx := jsonrequire.NewGetterWithJSON(a.t, string(a.details.RequestBody)).GetString("rawTx")

	tx, err := sdk.NewTransactionFromHex(rawTx)
	a.require.NoError(err)
	a.require.NotNil(tx)
	a.require.Equal(txID, tx.TxID().String())

	return a
}

func (a *arcBroadcastAssertions) WithCallbackURL(url string) ARCBroadcastAssertions {
	callbackURL, _ := a.details.RequestHeaders.Get("X-CallbackUrl")
	a.require.Equal(url, callbackURL)
	return a
}

func (a *arcBroadcastAssertions) WithCallbackToken(token string) ARCBroadcastAssertions {
	callbackToken, _ := a.details.RequestHeaders.Get("X-CallbackToken")
	a.require.Equal(token, callbackToken)
	return a
}
