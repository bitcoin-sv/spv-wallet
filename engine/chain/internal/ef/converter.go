package ef

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Converter provides a method to convert a transaction to EFHex format
type Converter struct {
	txsGetter chainmodels.TransactionsGetter
}

// NewConverter creates a new instance of Converter
func NewConverter(txsGetter chainmodels.TransactionsGetter) *Converter {
	return &Converter{txsGetter: txsGetter}
}

// Convert converts a (go-sdk) transaction to EFHex format
// Besides returning EFHex it also modifies the provided tx object to include missing inputs
// Missing source transactions for "unsourced" inputs are fetched using TransactionsGetter interface
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
