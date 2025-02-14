package testabilities

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type UTXOSelectorFixture interface {
	WillReturnNoUTXOs()
	WillReturnError()
}

var UserFundsTransactionOutpoint = bsv.Outpoint{
	TxID: "a0000000001e1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9",
	Vout: 0,
}

var UserFundsTransactionCustomInstructions = bsv.CustomInstructions{
	{Type: "type42", Instruction: "1-paymail_pki-0"},
	{Type: "type42", Instruction: "1-destination-0123"},
}

type mockedUTXOSelector struct {
	returnNothing bool
	returnError   bool
}

func (m *mockedUTXOSelector) Select(ctx context.Context, tx *sdk.Transaction, userID string) ([]*outlines.UTXO, error) {
	if m.returnError {
		return nil, spverrors.Newf("mocked: failed to select utxos for transaction")
	}

	if m.returnNothing {
		return nil, nil
	}

	fee := bsv.Satoshis(1)

	return []*outlines.UTXO{
		{
			TxID: UserFundsTransactionOutpoint.TxID,
			Vout: UserFundsTransactionOutpoint.Vout,
			CustomInstructions: bsv.CustomInstructions{
				{
					Type:        "type42",
					Instruction: "1-paymail_pki-0",
				},
				{
					Type:        "type42",
					Instruction: "1-destination-0123",
				},
			},
			Satoshis:           bsv.Satoshis(tx.TotalOutputSatoshis()) + fee,
			EstimatedInputSize: 148, // P2PKH input size
		},
	}, nil
}

func (m *mockedUTXOSelector) WillReturnNoUTXOs() {
	m.returnNothing = true
}

func (m *mockedUTXOSelector) WillReturnError() {
	m.returnError = true
}
