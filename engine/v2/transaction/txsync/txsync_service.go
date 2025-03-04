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

	trackedTx, err := s.transactionsRepo.GetTransaction(ctx, txInfo.TxID)
	if err != nil {
		return spverrors.Wrapf(err, "failed to get transaction %s", txInfo.TxID)
	}

	if trackedTx.UpdatedAt.After(txInfo.Timestamp) {
		s.logger.Info().Msgf("Tx %s has already been updated", txInfo.TxID)
		return nil
	}

	if txInfo.TXStatus.IsProblematic() {
		trackedTx.TxStatus = txmodels.TxStatusProblematic
		err = s.transactionsRepo.UpdateTransaction(ctx, trackedTx)
		if err != nil {
			return spverrors.Wrapf(err, "failed to set PROBLEMATIC status for transaction %s", txInfo.TxID)
		}
		return nil
	} else if !txInfo.TXStatus.IsMined() {
		s.logger.Info().
			Str("TxID", txInfo.TxID).
			Msgf("Received ARC callback with transaction which is not mined yet")
		return nil
	}

	bump, err := parseMerklePath(txInfo.MerklePath, txInfo.TxID)
	if err != nil {
		return err
	}

	if int64(bump.BlockHeight) != txInfo.BlockHeight {
		return spverrors.Newf("Block height in BUMP doesn't match the block height in the callback")
	}

	if trackedTx.BlockHash != nil && *trackedTx.BlockHash != txInfo.BlockHash {
		s.logger.Info().
			Str("TxID", txInfo.TxID).
			Msg("Received callback for already MINED transaction with different BUMP. Reorg could happen")
	}

	err = trackedTx.Mined(txInfo.BlockHash, bump)
	if err != nil {
		return spverrors.Wrapf(err, "failed to set MINED status for transaction %s", txInfo.TxID)
	}

	err = s.transactionsRepo.UpdateTransaction(ctx, trackedTx)
	if err != nil {
		return spverrors.Wrapf(err, "failed to set MINED status for transaction %s", txInfo.TxID)
	}

	return nil
}

func parseMerklePath(merklePath string, txID string) (*trx.MerklePath, error) {
	bump, err := trx.NewMerklePathFromHex(merklePath)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to parse merkle path for transaction %s", txID)
	}

	_, err = bump.ComputeRootHex(&txID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to validate merkle path %s for transaction %s", merklePath, txID)
	}

	return bump, nil
}
