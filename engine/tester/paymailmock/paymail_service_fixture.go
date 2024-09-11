package paymailmock

import (
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// PaymailClientServiceMock is a paymail.ServiceClient with mocked paymail.Client
type PaymailClientServiceMock struct {
	paymail.ServiceClient
	*PaymailClientMock
}

// CreatePaymailClientService creates a new paymail.ServiceClient with mocked paymail.Client
func CreatePaymailClientService(domain string, otherDomains ...string) *PaymailClientServiceMock {
	pmClient := MockClient(domain, otherDomains...)
	client, err := paymail.NewServiceClient(nil, pmClient)
	if err != nil {
		panic(spverrors.Wrapf(err, "cannot create paymail client service with mocked paymail"))
	}

	return &PaymailClientServiceMock{
		ServiceClient:     client,
		PaymailClientMock: pmClient,
	}
}
