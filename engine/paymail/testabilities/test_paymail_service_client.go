package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// PaymailClientGiven represents the operations that helps with prepare given state of the unit tests environment.
type PaymailClientGiven interface {
	NewPaymailClientService() paymailclient.ServiceClient
	MockedPaymailClient() *paymailmock.PaymailClientMock
	ExternalPaymailHost() PaymailHostGiven
}

// PaymailHostGiven represents the operations that helps with configuring the paymail host responses.
type PaymailHostGiven interface {
	MockedPaymailClient() *paymailmock.PaymailClientMock
	WillRespondWithP2PDestinationsWithSats(satoshis bsv.Satoshis, moreSatoshis ...bsv.Satoshis) *paymailmock.MockedP2PDestinationResponse
	WillRespondWithBasicCapabilities()
	WillRespondWithP2PCapabilities()
	WillRespondWithP2PWithBEEFCapabilities()
	WillRespondWithNotFoundOnCapabilities()
	WillRespondWithErrorOnCapabilities()
	WillRespondOnCapability(capabilityName string) *paymailmock.CapabilityMock
	WillRespondWithNotFoundOnP2PDestination()
	WillRespondWithErrorOnP2PDestinations()
}

type paymailServiceClientAbility struct {
	t *testing.T
	*paymailmock.PaymailClientMock
}

// New creates a new test ability.
func New(t testing.TB, domains ...string) (given PaymailClientGiven) {
	ability := &paymailServiceClientAbility{
		t:                 t.(*testing.T),
		PaymailClientMock: paymailmock.MockClient(fixtures.PaymailDomainExternal, domains...),
	}
	return ability
}

func (a *paymailServiceClientAbility) NewPaymailClientService() paymailclient.ServiceClient {
	return paymailclient.NewServiceClient(tester.CacheStore(), a.PaymailClientMock, tester.Logger(a.t))
}

func (a *paymailServiceClientAbility) MockedPaymailClient() *paymailmock.PaymailClientMock {
	return a.PaymailClientMock
}

func (a *paymailServiceClientAbility) ExternalPaymailHost() PaymailHostGiven {
	return a
}

func (a *paymailServiceClientAbility) WillRespondWithNotFoundOnP2PDestination() {
	a.PaymailClientMock.WillRespondWithP2PCapabilities()
	a.PaymailClientMock.WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).WithNotFound()
}

func (a *paymailServiceClientAbility) WillRespondWithErrorOnP2PDestinations() {
	a.PaymailClientMock.WillRespondWithP2PCapabilities()
	a.PaymailClientMock.WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).WithInternalServerError()
}

func (a *paymailServiceClientAbility) WillRespondWithP2PDestinationsWithSats(satoshis bsv.Satoshis, moreSatoshis ...bsv.Satoshis) *paymailmock.MockedP2PDestinationResponse {
	paymailHostResponse := paymailmock.P2PDestinationsForSats(satoshis, moreSatoshis...)

	a.PaymailClientMock.WillRespondWithP2PCapabilities()
	a.PaymailClientMock.
		WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
		With(paymailHostResponse)

	return paymailHostResponse
}
