package outlines

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/rs/zerolog"
)

// evaluationContext is a context for the evaluation of a transaction outline specification.
type evaluationContext interface {
	context.Context
	XPubID() string
	UserID() string
	Log() *zerolog.Logger
	Paymail() paymail.ServiceClient
	PaymailAddressService() paymailaddress.Service
}

type ctx struct {
	context.Context
	userID                string
	log                   *zerolog.Logger
	paymail               paymail.ServiceClient
	paymailAddressService paymailaddress.Service
}

// newTransactionContext creates a new context
func newTransactionContext(c context.Context, userID string, log *zerolog.Logger, paymail paymail.ServiceClient, paymailAddressService paymailaddress.Service) evaluationContext {
	return &ctx{
		Context:               c,
		userID:                userID,
		log:                   log,
		paymail:               paymail,
		paymailAddressService: paymailAddressService,
	}
}

func (c *ctx) XPubID() string {
	return ""
}

func (c *ctx) UserID() string {
	return c.userID
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
