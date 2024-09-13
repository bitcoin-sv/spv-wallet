package paymail

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
)

// ServiceClient is a service that aims to make easier paymail operations.
type ServiceClient interface {
	GetSanitizedPaymail(addr string) (*paymail.SanitisedPaymail, error)
	GetCapabilities(ctx context.Context, domain string) (*paymail.CapabilitiesPayload, error)
	GetP2P(ctx context.Context, domain string) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PayloadFormat)
	StartP2PTransaction(alias, domain, p2pDestinationURL string, satoshis uint64) (*paymail.PaymentDestinationPayload, error)
	GetPkiForPaymail(ctx context.Context, sPaymail *paymail.SanitisedPaymail) (*paymail.PKIResponse, error)
	AddContactRequest(ctx context.Context, receiverPaymail *paymail.SanitisedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error)
}
