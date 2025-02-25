package repository

import (
	"context"
	"iter"
	"slices"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Outputs is a repository for outputs.
type Outputs struct {
	db *gorm.DB
}

// NewOutputsRepo creates a new repository for outputs.
func NewOutputsRepo(db *gorm.DB) *Outputs {
	return &Outputs{db: db}
}

// FindByOutpoints returns outputs from the database based on the provided outpoints.
func (o *Outputs) FindByOutpoints(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]txmodels.TrackedOutput, error) {
	outpointsClause := slices.Collect(func(yield func(sqlPair []any) bool) {
		for outpoint := range outpoints {
			yield([]any{outpoint.TxID, outpoint.Vout})
		}
	})

	var outputs []*database.TrackedOutput

	if err := o.db.WithContext(ctx).
		Model(&database.TrackedOutput{}).
		Where("(tx_id, vout) IN ?", outpointsClause).
		Find(&outputs).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get outputs")
	}

	return lo.Map(outputs, func(output *database.TrackedOutput, _ int) txmodels.TrackedOutput {
		return txmodels.TrackedOutput{
			TxID: output.TxID,
			Vout: output.Vout,
			// Because the field is char(64) in the database, it will be padded with spaces.
			// For some reason the value is padded only on postgres,
			// so if changing, make sure to check it on all databases.
			SpendingTX: strings.TrimLeft(output.SpendingTX, " "),
			UserID:     output.UserID,
			Satoshis:   output.Satoshis,
			CreatedAt:  output.CreatedAt,
			UpdatedAt:  output.UpdatedAt,
		}
	}), nil
}
