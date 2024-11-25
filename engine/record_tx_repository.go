package engine

import (
	"context"
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type recordTXRepository struct {
	db *gorm.DB
}

// SaveTX saves a transaction to the database.
func (r *recordTXRepository) SaveTX(ctx context.Context, txRow *database.Transaction) error {
	query := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		})

	if err := query.Create(txRow).Error; err != nil {
		return spverrors.Wrapf(err, "failed to save transaction")
	}

	return nil
}

// GetOutputs returns outputs from the database based on the provided outpoints.
func (r *recordTXRepository) GetOutputs(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error) {
	outpointsClause := slices.Collect(func(yield func(sqlPair []any) bool) {
		for outpoint := range outpoints {
			yield([]any{outpoint.TxID, outpoint.Vout})
		}
	})

	query := r.db.
		WithContext(ctx).
		Model(&database.Output{}).
		Where("(tx_id, vout) IN ?", outpointsClause)

	var outputs []*database.Output
	if err := query.Find(&outputs).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get outputs")
	}

	return outputs, nil
}
