package internal

import (
	"context"
	"maps"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
)

// CombinedTxsGetter is a TransactionsGetter that combines multiple TransactionsGetters
type CombinedTxsGetter struct {
	txsGetters []chainmodels.TransactionsGetter
}

// NewCombinedTxsGetter creates a new CombinedTxsGetter
func NewCombinedTxsGetter(txsGetters ...chainmodels.TransactionsGetter) *CombinedTxsGetter {
	return &CombinedTxsGetter{
		txsGetters: txsGetters,
	}
}

// GetTransactions gets transactions from all provided TransactionsGetters in order
// the first tx getter is queried for all transactions, the second tx getter is queried only for the missing transactions and so on
func (ctg *CombinedTxsGetter) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	missingTxs := map[string]bool{}
	for id := range ids {
		missingTxs[id] = true
	}
	var transactions []*sdk.Transaction
	for _, getter := range ctg.txsGetters {
		if len(missingTxs) == 0 {
			break
		}
		if getter == nil {
			return nil, spverrors.Newf("nil transactions getter")
		}
		txs, err := getter.GetTransactions(ctx, maps.Keys(missingTxs))
		if err != nil {
			return nil, spverrors.ErrGetTransactions.Wrap(err)
		}
		for _, tx := range txs {
			txID := tx.TxID().String()
			if _, exists := missingTxs[txID]; !exists {
				// This transaction was already fetched by another getter
				continue
			}
			delete(missingTxs, txID)
			transactions = append(transactions, tx)
		}
	}
	return transactions, nil
}
