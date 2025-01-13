package repository

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iter"
	"slices"
)

type Operations struct {
	db *gorm.DB
}

func NewOperationsRepo(db *gorm.DB) *Operations {
	return &Operations{db: db}
}

func (o *Operations) PaginatedForUser(ctx context.Context, userID string, page filter.Page) (*database.PagedResult[database.Operation], error) {
	return database.PaginatedQuery[database.Operation](ctx, page, o.db, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("user_id = ?", userID)
	})
}

// SaveAll saves operations to the database.
func (o *Operations) SaveAll(ctx context.Context, opRows iter.Seq[*database.Operation]) error {
	query := o.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		})

	if err := query.Create(slices.Collect(opRows)).Error; err != nil {
		return err
	}

	return nil
}
