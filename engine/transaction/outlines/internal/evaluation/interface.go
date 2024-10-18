package evaluation

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/rs/zerolog"
)

// Context is a context for the evaluation of a transaction outline specification.
type Context interface {
	context.Context
	XPubID() string
	Log() *zerolog.Logger
	Paymail() paymail.ServiceClient
	PaymailAddressService() paymailaddress.Service
}
