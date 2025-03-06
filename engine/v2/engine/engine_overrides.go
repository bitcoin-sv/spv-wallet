package engine

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/go-resty/resty/v2"
)

// InternalsOverride is a function that can be used to override internal dependencies.
// This is meant to be used for testing purposes.
type InternalsOverride = func(*overrides)

type overrides struct {
	transport     http.RoundTripper
	resty         *resty.Client
	paymailClient paymail.ClientInterface
}

// WithResty is a function that can be used to override the resty.Client used by the engine.
// This is meant to be used for testing purposes.
func WithResty(resty *resty.Client) InternalsOverride {
	return func(o *overrides) {
		o.resty = resty
	}
}

// WithTransport is a function that can be used to override the http.RoundTripper used by the engine.
// This is meant to be used for testing purposes.
func WithTransport(transport http.RoundTripper) InternalsOverride {
	return func(o *overrides) {
		o.transport = transport
	}
}

// WithPaymailClient is a function that can be used to override the paymail.ClientInterface used by the engine.
// This is meant to be used for testing purposes.
func WithPaymailClient(client paymail.ClientInterface) InternalsOverride {
	return func(o *overrides) {
		o.paymailClient = client
	}
}
