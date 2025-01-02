package dao

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iter"
	"slices"
)

type Operations struct {
	db *gorm.DB
}

func NewOperationsAccessObject(db *gorm.DB) *Operations {
	return &Operations{db: db}
}

func (r *Operations) SaveOperation(ctx context.Context, opRows iter.Seq[*database.Operation]) error {
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
