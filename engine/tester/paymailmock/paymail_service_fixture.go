package paymailmock

import (
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
)

// PaymailClientServiceMock is a paymail.ServiceClient with mocked paymail.Client
type PaymailClientServiceMock struct {
	paymail.ServiceClient
	*PaymailClientMock
}

// CreatePaymailClientService creates a new paymail.ServiceClient with mocked paymail.Client
func CreatePaymailClientService(domain string, otherDomains ...string) *PaymailClientServiceMock {
	pmClient := MockClient(domain, otherDomains...)
	client := paymail.NewServiceClient(nil, pmClient)

	return &PaymailClientServiceMock{
		ServiceClient:     client,
		PaymailClientMock: pmClient,
	}
}
