package record

import (
	"iter"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

type customOutputsResolver struct {
	flow        *txFlow
	userID      string
	userPubKey  *primitives.PublicKey
	annotations transaction.OutputsAnnotations
	txAddresses map[string]uint32
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

func (c *customOutputsResolver) remainingAddresses() map[string]uint32 {
	return c.txAddresses
}

func (c *customOutputsResolver) customOutputs() iter.Seq[txmodels.NewOutput] {
	return func(yield func(output txmodels.NewOutput) bool) {
		for vout, annotation := range c.annotations {
			if annotation.Bucket != bucket.BSV || annotation.CustomInstructions == nil {
				continue
			}

			userPubKey, err := c.getUserPubKey()
			if err != nil {
				c.err = spverrors.Wrapf(err, "failed to get public key for user %s", c.userID)
				break
			}
			calculatedAddr, err := custominstructions.Address(*userPubKey, *annotation.CustomInstructions)
			if err != nil {
				c.err = spverrors.Wrapf(err, "failed to derive address from custom instructions for user %s", c.userID)
				break
			}

			realVOut, ok := c.txAddresses[calculatedAddr.AddressString]
			if !ok {
				c.err = spverrors.Newf("address derived from custom instructions doesn't match any of the addresses in the locking scripts")
				break
			}

			if realVOut != vout {
				c.err = spverrors.Newf("address derived from custom instructions doesn't match the address in the locking script")
				break
			}

			yield(txmodels.NewOutputForP2PKH(
				bsv.Outpoint{TxID: c.flow.txID, Vout: realVOut},
				c.userID,
				bsv.Satoshis(c.flow.tx.Outputs[vout].Satoshis),
				*annotation.CustomInstructions,
			))

			delete(c.txAddresses, calculatedAddr.AddressString)
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
