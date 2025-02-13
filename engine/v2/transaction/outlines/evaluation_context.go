package outlines

import (
	"context"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"

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
	feeUnit               bsvmodel.FeeUnit
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
