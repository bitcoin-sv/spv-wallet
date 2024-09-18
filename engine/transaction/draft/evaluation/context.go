package evaluation

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/rs/zerolog"
)

type ctx struct {
	context.Context
	log     *zerolog.Logger
	paymail paymail.ServiceClient
}

// NewContext creates a new context
func NewContext(c context.Context, log *zerolog.Logger, paymail paymail.ServiceClient) Context {
	return &ctx{Context: c, log: log, paymail: paymail}
}

func (c *ctx) Log() *zerolog.Logger {
	return c.log
}

func (c *ctx) Paymail() paymail.ServiceClient {
	return c.paymail
}
