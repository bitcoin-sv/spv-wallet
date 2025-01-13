package repository

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
	"iter"
	"slices"
)

type Outputs struct {
	db *gorm.DB
}

func NewOutputsRepo(db *gorm.DB) *Outputs {
	return &Outputs{db: db}
}

// FindByOutpoints returns outputs from the database based on the provided outpoints.
func (o *Outputs) FindByOutpoints(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.UserUtxos, []*database.TrackedOutput, error) {
	outpointsClause := slices.Collect(func(yield func(sqlPair []any) bool) {
		for outpoint := range outpoints {
			yield([]any{outpoint.TxID, outpoint.Vout})
		}
	})

	var utxos []*database.UserUtxos
	var trackedOutputs []*database.TrackedOutput
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&database.UserUtxos{}).
			Where("(tx_id, vout) IN ?", outpointsClause).
			Find(&utxos).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&database.TrackedOutput{}).
			Where("(tx_id, vout) IN ?", outpointsClause).
			Find(&trackedOutputs).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, spverrors.Wrapf(err, "failed to get outputs")
	}

	return utxos, trackedOutputs, nil
}
