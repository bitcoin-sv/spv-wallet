package internal

import (
	"context"
	"iter"
	"maps"
	"slices"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// CombineTxsGetters creates a new CombinedTxsGetter
func CombineTxsGetters(txsGetters ...chainmodels.TransactionsGetter) chainmodels.TransactionsGetter {
	getters := filterNilGetters(txsGetters...)

	if len(getters) == 0 {
		return &emptyGetter{}
	}
	if len(getters) == 1 {
		return getters[0]
	}
	return &combinedTxsGetter{
		txsGetters: getters,
	}
}

type combinedTxsGetter struct {
	txsGetters []chainmodels.TransactionsGetter
}

// GetTransactions gets transactions from all provided TransactionsGetters in order
// the first tx getter is queried for all transactions, the second tx getter is queried only for the missing transactions and so on
func (ctg *combinedTxsGetter) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
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
			return nil, chainerrors.ErrGetTransactionsByTxsGetter.Wrap(err)
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

type emptyGetter struct{}

func (ctg *emptyGetter) GetTransactions(_ context.Context, _ iter.Seq[string]) ([]*sdk.Transaction, error) {
	return nil, nil
}

func filterNilGetters(txsGetters ...chainmodels.TransactionsGetter) []chainmodels.TransactionsGetter {
	return slices.DeleteFunc(txsGetters, func(getter chainmodels.TransactionsGetter) bool {
		return getter == nil
	})
}
