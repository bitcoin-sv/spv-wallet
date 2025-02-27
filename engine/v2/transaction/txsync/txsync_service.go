package txsync

import (
	"context"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

type Service struct {
	transactionsRepo TransactionsRepo
}

func NewService(transactionsRepo TransactionsRepo) *Service {
	return &Service{
		transactionsRepo: transactionsRepo,
	}
}

func (s *Service) Handle(ctx context.Context, txInfo chainmodels.TXInfo) error {
	if txInfo.TXStatus.IsProblematic() {
		err := s.transactionsRepo.SetStatus(ctx, txInfo.TxID, txmodels.TxStatusProblematic)
		if err != nil {
			return spverrors.Wrapf(err, "failed to set PROBLEMATIC status for transaction %s", txInfo.TxID)
		}
	}

	return nil
}
