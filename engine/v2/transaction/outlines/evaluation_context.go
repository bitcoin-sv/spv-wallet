package outlines

import (
	"context"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"
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
	usersService          UsersService
}

func (c *evaluationContext) UserID() string {
	return c.userID
}

func (c *evaluationContext) UserPubKey() (*primitives.PublicKey, error) {
	pubKey, err := c.usersService.GetPubKey(c, c.userID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get public key for user %s", c.userID)
	}
	return pubKey, nil
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
