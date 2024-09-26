package evaluation

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/rs/zerolog"
)

type ctx struct {
	context.Context
	xPubID                string
	log                   *zerolog.Logger
	paymail               paymail.ServiceClient
	paymailAddressService paymailaddress.Service
}

// NewContext creates a new context
func NewContext(c context.Context, xPubID string, log *zerolog.Logger, paymail paymail.ServiceClient, paymailAddressService paymailaddress.Service) Context {
	return &ctx{
		Context:               c,
		xPubID:                xPubID,
		log:                   log,
		paymail:               paymail,
		paymailAddressService: paymailAddressService,
	}
}

func (c *ctx) XPubID() string {
	return c.xPubID
}

func (c *ctx) Log() *zerolog.Logger {
	return c.log
}

func (c *ctx) Paymail() paymail.ServiceClient {
	return c.paymail
}

func (c *ctx) PaymailAddressService() paymailaddress.Service {
	return c.paymailAddressService
}
