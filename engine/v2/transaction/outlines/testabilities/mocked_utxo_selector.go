package testabilities

import (
	"context"
	"fmt"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
)

type UTXOSelectorFixture interface {
	WillReturnNoUTXOs()
	WillReturnError()
	WillReturnUTXOs(change bsv.Satoshis, utxos ...bsv.Satoshis)
}

func templatedOutpoint(index uint) bsv.Outpoint {
	return bsv.Outpoint{
		TxID: fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", index),
		Vout: 0,
	}
}

var UserFundsTransactionOutpoint = templatedOutpoint(0)

var UserFundsTransactionCustomInstructions = bsv.CustomInstructions{
	{Type: "type42", Instruction: "1-paymail_pki-0"},
	{Type: "type42", Instruction: "1-destination-0123"},
}

type mockedUTXOSelector struct {
	returnNothing  bool
	returnError    bool
	utxosToReturn  []bsv.Satoshis
	changeToReturn bsv.Satoshis
}

func (m *mockedUTXOSelector) Select(ctx context.Context, tx *sdk.Transaction, userID string) ([]*outlines.UTXO, bsv.Satoshis, error) {
	if m.returnError {
		return nil, 0, spverrors.Newf("mocked: failed to select utxos for transaction")
	}

	if m.returnNothing {
		return nil, 0, nil
	}

	var distribution []bsv.Satoshis
	if m.utxosToReturn != nil {
		distribution = m.utxosToReturn
	} else {
		// default case, produce no change
		fee := bsv.Satoshis(1)
		distribution = []bsv.Satoshis{
			bsv.Satoshis(tx.TotalOutputSatoshis()) + fee,
		}
	}

	return lo.Map(distribution, func(satoshis bsv.Satoshis, index int) *outlines.UTXO {
		outpoint := templatedOutpoint(uint(index))
		return &outlines.UTXO{
			TxID:               outpoint.TxID,
			Vout:               outpoint.Vout,
			CustomInstructions: UserFundsTransactionCustomInstructions,
			Satoshis:           satoshis,
			EstimatedInputSize: 148, // P2PKH input size
		}
	}), m.changeToReturn, nil
}

func (m *mockedUTXOSelector) WillReturnNoUTXOs() {
	m.returnNothing = true
}

func (m *mockedUTXOSelector) WillReturnError() {
	m.returnError = true
}

func (m *mockedUTXOSelector) WillReturnUTXOs(change bsv.Satoshis, utxos ...bsv.Satoshis) {
	m.utxosToReturn = utxos
	m.changeToReturn = change
}
