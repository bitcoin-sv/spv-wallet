package evaluation

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/rs/zerolog"
)

// Context is a context for the evaluation of a transaction draft specification.
type Context interface {
	context.Context
	Log() *zerolog.Logger
	Paymail() paymail.ServiceClient
}
