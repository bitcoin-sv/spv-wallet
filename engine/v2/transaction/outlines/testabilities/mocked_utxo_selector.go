package testabilities

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"
)

var UserFundsTransactionOutpoint = bsvmodel.Outpoint{
	TxID: "a0000000001e1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9",
	Vout: 0,
}

type mockedUTXOSelector struct {
	returnNothing bool
}

func (m *mockedUTXOSelector) Select(ctx context.Context, tx *sdk.Transaction, userID string) ([]*bsvmodel.Outpoint, error) {
	if m.returnNothing {
		return nil, nil
	}

	return []*bsvmodel.Outpoint{
		{
			TxID: "a0000000001e1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9",
			Vout: 0,
		},
	}, nil
}

func (m *mockedUTXOSelector) WillReturnNoUTXOs() {
	m.returnNothing = true
}
