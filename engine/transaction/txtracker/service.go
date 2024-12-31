package txtracker

import (
	"context"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"gorm.io/datatypes"
	"iter"
	"maps"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) TrackMissingTxs(ctx context.Context, transactions iter.Seq[*trx.Transaction]) error {
	txsMap := maps.Collect(func(yield func(string, *trx.Transaction) bool) {
		for tx := range transactions {
			yield(tx.TxID().String(), tx)
		}
	})

	missingTxIDs, err := s.repo.MissingTransactions(ctx, maps.Keys(txsMap))
	if err != nil {
		return ErrCannotCheckMissingTransactions.Wrap(err)
	}

	transactionModels := func(yield func(transaction *database.TrackedTransaction) bool) {
		for txID := range missingTxIDs {
			tx := txsMap[txID]
			isMined := tx.MerklePath != nil

			row := &database.TrackedTransaction{
				ID:       txID,
				TxStatus: database.TxStatusCreated,
				//TODO: Question: What if tx is not mined? Should we try to rebroadcast it? Or create another Status, e.g. "tracked"
			}

			if isMined {
				row.TxStatus = database.TxStatusMined
				bumpJSONData := datatypes.NewJSONType(*tx.MerklePath)
				row.BUMP = &bumpJSONData
			}

			yield(row)
		}
	}

	err = s.repo.SaveTXs(ctx, transactionModels)
	if err != nil {
		return ErrCannotSaveTransactions.Wrap(err)
	}

	return nil
}
