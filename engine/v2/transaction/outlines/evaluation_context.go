package outlines

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/rs/zerolog"
)

type evaluationContext struct {
	context.Context
	userID                string
	log                   *zerolog.Logger
	paymail               paymail.ServiceClient
	paymailAddressService PaymailAddressService
	utxoSelector          UTXOSelector
}

func newOutlineEvaluationContext(ctx context.Context, userID string, log *zerolog.Logger, paymail paymail.ServiceClient, paymailAddressService PaymailAddressService, utxoSelector UTXOSelector) *evaluationContext {
	return &evaluationContext{
		Context:               ctx,
		userID:                userID,
		log:                   log,
		paymail:               paymail,
		paymailAddressService: paymailAddressService,
		utxoSelector:          utxoSelector,
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

func (c *evaluationContext) PaymailAddressService() PaymailAddressService {
	return c.paymailAddressService
}

func (c *evaluationContext) UTXOSelector() UTXOSelector {
	return c.utxoSelector
}
