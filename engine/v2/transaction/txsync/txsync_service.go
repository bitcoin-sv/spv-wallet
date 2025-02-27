package txsync

import (
	"context"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/rs/zerolog"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

type Service struct {
	logger           zerolog.Logger
	transactionsRepo TransactionsRepo
}

func NewService(logger zerolog.Logger, transactionsRepo TransactionsRepo) *Service {
	return &Service{
		transactionsRepo: transactionsRepo,
		logger:           logger,
	}
}

func (s *Service) Handle(ctx context.Context, txInfo chainmodels.TXInfo) error {
	if txInfo.TxID == "" {
		return spverrors.Newf("Received ARC callback with empty transaction ID")
	}

	if txInfo.TXStatus.IsProblematic() {
		err := s.transactionsRepo.SetStatus(ctx, txInfo.TxID, txmodels.TxStatusProblematic)
		if err != nil {
			return spverrors.Wrapf(err, "failed to set PROBLEMATIC status for transaction %s", txInfo.TxID)
		}
		return nil
	}

	if !txInfo.TXStatus.IsMined() {
		s.logger.Info().
			Str("TxID", txInfo.TxID).
			Msgf("Received ARC callback with transaction which is not mined yet")
	}

	bump, err := trx.NewMerklePathFromHex(txInfo.MerklePath)
	if err != nil {
		return spverrors.Wrapf(err, "failed to parse merkle path for transaction %s", txInfo.TxID)
	}
	_ = bump

	err = s.transactionsRepo.SetAsMined(ctx, txInfo.TxID, txInfo.BlockHash, txInfo.BlockHeight)
	if err != nil {
		return spverrors.Wrapf(err, "failed to set MINED status for transaction %s", txInfo.TxID)
	}

	return nil
}
