package repository

import (
	"context"
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/dbquery"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Operations is a repository for operations.
type Operations struct {
	db *gorm.DB
}

// NewOperationsRepo creates a new repository for operations.
func NewOperationsRepo(db *gorm.DB) *Operations {
	return &Operations{db: db}
}

// PaginatedForUser returns operations for a user based on userID and the provided paging options.
func (o *Operations) PaginatedForUser(ctx context.Context, userID string, page filter.Page) (*dbquery.PagedResult[database.Operation], error) {
	return dbquery.PaginatedQuery[database.Operation](ctx, page, o.db, dbquery.UserID(userID))
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
