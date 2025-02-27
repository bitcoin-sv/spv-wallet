package txsync

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/rs/zerolog"
)

// Service is meant to handle the ARC callback and update the transaction status in the database.
type Service struct {
	logger           zerolog.Logger
	transactionsRepo TransactionsRepo
}

// NewService creates a new transaction sync service.
func NewService(logger zerolog.Logger, transactionsRepo TransactionsRepo) *Service {
	return &Service{
		transactionsRepo: transactionsRepo,
		logger:           logger,
	}
}

// Handle processes the ARC callback and updates the transaction status in the database.
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

	hex, isBEEF, err := s.transactionsRepo.GetTransactionHex(ctx, txInfo.TxID)
	if err != nil {
		return spverrors.Wrapf(err, "failed to get transaction hex for transaction %s", txInfo.TxID)
	}
	var tx *trx.Transaction
	if isBEEF {
		tx, err = trx.NewTransactionFromBEEFHex(hex)
	} else {
		tx, err = trx.NewTransactionFromHex(hex)
	}
	if err != nil {
		return spverrors.Wrapf(err, "failed to parse transaction hex for transaction %s", txInfo.TxID)
	}

	tx.MerklePath = bump

	beefHex, err := tx.BEEFHex()
	if err != nil {
		return spverrors.Wrapf(err, "failed to get BEEF hex for transaction %s", txInfo.TxID)
	}

	err = s.transactionsRepo.SetAsMined(ctx, txInfo.TxID, txInfo.BlockHash, txInfo.BlockHeight, beefHex)
	if err != nil {
		return spverrors.Wrapf(err, "failed to set MINED status for transaction %s", txInfo.TxID)
	}

	return nil
}
