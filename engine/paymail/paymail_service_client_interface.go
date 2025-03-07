package paymail

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// Cache is an interface that defines the methods to interact with a cache.
type Cache interface {
	GetModel(ctx context.Context, key string, model interface{}) error
	SetModel(ctx context.Context, key string, model interface{}, ttl time.Duration, dependencies ...string) error
}

// ServiceClient is a service that aims to make easier paymail operations.
type ServiceClient interface {
	GetSanitizedPaymail(addr string) (*paymail.SanitisedPaymail, error)
	GetCapabilities(ctx context.Context, domain string) (*paymail.CapabilitiesPayload, error)
	GetP2PDestinations(ctx context.Context, address *paymail.SanitisedPaymail, satoshis bsv.Satoshis) (*paymail.PaymentDestinationPayload, error)
	GetP2P(ctx context.Context, domain string) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PayloadFormat)
	StartP2PTransaction(alias, domain, p2pDestinationURL string, satoshis uint64) (*paymail.PaymentDestinationPayload, error)
	GetPkiForPaymail(ctx context.Context, sPaymail *paymail.SanitisedPaymail) (*paymail.PKIResponse, error)
	AddContactRequest(ctx context.Context, receiverPaymail *paymail.SanitisedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error)
	Notify(ctx context.Context, address string, p2pMetadata *paymail.P2PMetaData, reference string, tx *trx.Transaction) error
}
