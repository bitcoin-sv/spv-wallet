package ef

import (
	"context"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
)

type TransactionsGetter interface {
	GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error)
}

type Converter struct {
	txsGetter TransactionsGetter
}

func NewConverter(txsGetter TransactionsGetter) *Converter {
	return &Converter{txsGetter: txsGetter}
}

func (c *Converter) Convert(ctx context.Context, tx *sdk.Transaction) (string, error) {
	/**
	NOTE: We can't make out-of-the-box check if the tx is already in EF format (because sourceOutput field of tx object is unexported)
	but we can try to convert to EFHex and return it if it's possible,
	otherwise we will try to find the missing inputs and hydrate them
	*/

	efHex, err := tx.EFHex()
	if err == nil {
		return efHex, nil
	}

	unsourcedInputs, err := findUnsourcedInputs(tx)
	if err != nil {
		return "", err
	}

	sourceTransactions, err := c.txsGetter.GetTransactions(ctx, unsourcedInputs.getMissingTXIDs())
	if err != nil {
		return "", ErrGetTransactions.Wrap(err)
	}
	if len(sourceTransactions) != unsourcedInputs.txCount() {
		return "", ErrGetTransactions.Wrap(spverrors.Newf("expected %d transactions, got %d", unsourcedInputs.txCount(), len(sourceTransactions)))
	}

	for _, sourceTX := range sourceTransactions {
		if sourceTX == nil {
			return "", ErrGetTransactions.Wrap(spverrors.Newf("nil transaction found"))
		}
		if err := unsourcedInputs.hydrate(sourceTX); err != nil {
			return "", ErrGetTransactions.Wrap(err)
		}
		unsourcedInputs.deleteTXID(sourceTX.TxID().String())
	}

	return makeEFHex(tx)
}

func makeEFHex(tx *sdk.Transaction) (string, error) {
	efHex, err := tx.EFHex()
	if err != nil {
		return "", ErrEFHexGeneration.Wrap(err)
	}
	return efHex, nil
}
