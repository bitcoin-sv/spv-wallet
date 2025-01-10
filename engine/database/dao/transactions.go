package dao

import (
	"context"
	"iter"
	"maps"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Transactions is a data access object for transactions.
type Transactions struct {
	db *gorm.DB
}

// NewTransactionsAccessObject creates a new access object for transactions.
func NewTransactionsAccessObject(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

// SaveTXs saves transactions to the database.
func (r *Transactions) SaveTXs(ctx context.Context, txRows iter.Seq[*database.TrackedTransaction]) error {
	query := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		})

	if err := query.Create(slices.Collect(txRows)).Error; err != nil {
		return spverrors.Wrapf(err, "failed to save transaction")
	}

	return nil
}

// GetOutputs returns outputs from the database based on the provided outpoints.
func (r *Transactions) GetOutputs(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.UserUtxos, []*database.TrackedOutput, error) {
	outpointsClause := slices.Collect(func(yield func(sqlPair []any) bool) {
		for outpoint := range outpoints {
			yield([]any{outpoint.TxID, outpoint.Vout})
		}
	})

	var utxos []*database.UserUtxos
	var trackedOutputs []*database.TrackedOutput
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

// GetAddresses returns address entities from the database based on the provided address iterator.
func (r *Transactions) GetAddresses(ctx context.Context, addresses iter.Seq[string]) ([]*database.Address, error) {
	var rows []*database.Address
	if err := r.db.
		WithContext(ctx).
		Model(&database.Address{}).
		Where("address IN ?", slices.Collect(addresses)).
		Find(&rows).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get addresses")
	}

	return rows, nil
}

// MissingTransactions returns transactions that are not tracked in the database.
func (r *Transactions) MissingTransactions(ctx context.Context, txIDs iter.Seq[string]) (iter.Seq[string], error) {
	idsMap := maps.Collect(func(yield func(string, bool) bool) {
		for txID := range txIDs {
			yield(txID, true)
		}
	})

	idsSlice := slices.Collect(maps.Keys(idsMap))

	var alreadyTracked []string
	if err := r.db.
		WithContext(ctx).
		Model(&database.TrackedTransaction{}).
		Where("id IN ?", idsSlice).
		Pluck("id", &alreadyTracked).
		Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get missing transactions")
	}

	for _, txID := range alreadyTracked {
		delete(idsMap, txID)
	}

	return maps.Keys(idsMap), nil
}

// SaveOperations saves operations to the database.
func (r *Transactions) SaveOperations(ctx context.Context, opRows iter.Seq[*database.Operation]) error {
	query := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		})

	if err := query.Create(slices.Collect(opRows)).Error; err != nil {
		return err
	}

	return nil
}
