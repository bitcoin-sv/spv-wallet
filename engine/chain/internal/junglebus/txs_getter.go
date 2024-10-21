package junglebus

import (
	"context"
	"errors"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	chainerrors "github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
)

// GetTransactions implements chainmodels.TransactionsGetter interface to allow fetching transactions from Junglebus
func (s *Service) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	var transactions []*sdk.Transaction
	for id := range ids {
		select {
		case <-ctx.Done():
			return nil, spverrors.ErrCtxInterrupted.Wrap(ctx.Err())
		default:
			tx, err := s.FetchTransaction(ctx, id)
			if errors.Is(err, chainerrors.ErrJunglebusTxNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, tx)
		}
	}
	return transactions, nil
}
