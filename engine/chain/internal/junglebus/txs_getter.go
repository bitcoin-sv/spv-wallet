package junglebus

import (
	"context"
	"errors"
	chainerrors "github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"iter"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

func (s *Service) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	var transactions []*sdk.Transaction
	for id := range ids {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
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
