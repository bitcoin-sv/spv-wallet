package paymailmock

import (
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	xtester "github.com/bitcoin-sv/spv-wallet/engine/tester"
)

// PaymailClientServiceMock is a paymail.ServiceClient with mocked paymail.Client
type PaymailClientServiceMock struct {
	paymail.ServiceClient
	*PaymailClientMock
}

// CreatePaymailClientService creates a new paymail.ServiceClient with mocked paymail.Client
func CreatePaymailClientService(domain string, otherDomains ...string) *PaymailClientServiceMock {
	pmClient := MockClient(domain, otherDomains...)
	client := paymail.NewServiceClient(nil, pmClient, xtester.Logger())

	return &PaymailClientServiceMock{
		ServiceClient:     client,
		PaymailClientMock: pmClient,
	}
}
