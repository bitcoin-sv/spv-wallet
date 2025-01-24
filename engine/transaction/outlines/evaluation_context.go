package outlines

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/rs/zerolog"
)

type evaluationContext struct {
	context.Context
	userID                string
	log                   *zerolog.Logger
	paymail               paymail.ServiceClient
	paymailAddressService paymailaddress.Service
}

func newOutlineEvaluationContext(ctx context.Context, userID string, log *zerolog.Logger, paymail paymail.ServiceClient, paymailAddressService paymailaddress.Service) *evaluationContext {
	return &evaluationContext{
		Context:               ctx,
		userID:                userID,
		log:                   log,
		paymail:               paymail,
		paymailAddressService: paymailAddressService,
	}
}

func (c *evaluationContext) UserID() string {
	return c.userID
}

func (c *evaluationContext) Log() *zerolog.Logger {
	return c.log
}

func (c *evaluationContext) Paymail() paymail.ServiceClient {
	return c.paymail
}

func (c *evaluationContext) PaymailAddressService() paymailaddress.Service {
	return c.paymailAddressService
}
