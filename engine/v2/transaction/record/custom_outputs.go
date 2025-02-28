package record

import (
	"iter"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type customOutputsResolver struct {
	flow        *txFlow
	userID      string
	userPubKey  *primitives.PublicKey
	annotations transaction.OutputsAnnotations
	txAddresses addresses
	err         error
}

func (c *customOutputsResolver) getUserPubKey() (*primitives.PublicKey, error) {
	if c.userPubKey == nil {
		pubKey, err := c.flow.service.users.GetPubKey(c.flow.ctx, c.userID)
		if err != nil {
			return nil, spverrors.Wrapf(err, "failed to get public key for user %s", c.userID)
		}
		c.userPubKey = pubKey
	}
	return c.userPubKey, nil
}

func (c *customOutputsResolver) remainingAddresses() addresses {
	return c.txAddresses
}

func (c *customOutputsResolver) resolveAddress(address string, vout uint32) {
	c.txAddresses.remove(address, vout)
}

func (c *customOutputsResolver) annotatedOutputs() iter.Seq2[string, txmodels.NewOutput] {
	return func(yield func(address string, output txmodels.NewOutput) bool) {
		for vout, annotation := range c.annotations {
			if annotation.CustomInstructions == nil {
				continue
			}

			userPubKey, err := c.getUserPubKey()
			if err != nil {
				c.err = spverrors.Wrapf(err, "failed to get public key for user %s", c.userID)
				break
			}

			interpreted, err := custominstructions.NewAddressInterpreter().
				Process(userPubKey, *annotation.CustomInstructions)
			if err != nil {
				c.err = spverrors.Wrapf(err, "failed to derive address from custom instructions for user %s", c.userID)
				break
			}

			if ok := c.txAddresses.contains(interpreted.Address.AddressString, vout); !ok {
				c.err = spverrors.Newf("address derived from custom instructions doesn't match the locking script of the output")
				break
			}

			yield(interpreted.Address.AddressString, txmodels.NewOutputForP2PKH(
				bsv.Outpoint{TxID: c.flow.txID, Vout: vout},
				c.userID,
				bsv.Satoshis(c.flow.tx.Outputs[vout].Satoshis),
				*annotation.CustomInstructions,
			))
		}
	}
}

func (f *txFlow) resolveCustomOutputs(userID string, annotations transaction.OutputsAnnotations) *customOutputsResolver {
	return &customOutputsResolver{
		flow:        f,
		userID:      userID,
		txAddresses: f.allP2PKHAddresses(),
		annotations: annotations,
	}
}
