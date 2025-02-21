package repository

import (
	"context"
	"maps"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"gorm.io/gorm"
)

// Transactions provides database operations for managing transactions.
type Transactions struct {
	db *gorm.DB
}

// NewTransactions creates a new Transactions repository.
// It initializes a database-backed service for querying and managing transactions.
func NewTransactions(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

// HasTransactionInputSources checks if all the provided input source transaction IDs exist in the database.
// If all of them are found, the transaction data can be serialized into Raw HEX format.
// Otherwise, serialization should be done using the BEEFHex format.
func (t *Transactions) HasTransactionInputSources(ctx context.Context, sourceTXIDs ...string) (bool, error) {
	var count int64

	uniqueIds := stringSlice(sourceTXIDs).unique()
	tx := t.db.Session(&gorm.Session{PrepareStmt: true})
	err := tx.
		Model(&database.TrackedTransaction{}).
		WithContext(ctx).
		Where("id IN (?)", uniqueIds).
		Count(&count).Error
	if err != nil {
		return false, spverrors.Wrapf(err, "database query failed for source transactions %v", sourceTXIDs)
	}

	return count == int64(len(uniqueIds)), nil
}

// FindTransactionInputSources retrieves the full ancestry of input sources for a given transaction.
// It recursively traces input sources in batches to minimize database queries.
func (t *Transactions) FindTransactionInputSources(ctx context.Context, sourceTXIDs ...string) (beef.TxQueryResultSlice, error) {
	if len(sourceTXIDs) == 0 {
		return nil, nil
	}

	var transactions []database.TrackedTransaction
	uniqueIds := stringSlice(sourceTXIDs).unique()
	tx := t.db.Session(&gorm.Session{PrepareStmt: true})
	err := tx.
		WithContext(ctx).
		Raw(`
			WITH RECURSIVE transaction_tree AS (
				-- Start with the initial transactions
				SELECT id, beef_hex, raw_hex
				FROM xapi_tracked_transactions
				WHERE id IN (?)

				UNION ALL

				-- Recursive part: Traverse through source transactions
				SELECT st.id, st.beef_hex, st.raw_hex
				FROM xapi_tracked_transactions st
				INNER JOIN xapi_source_transactions st_map ON st.id = st_map.source_transaction_id
				INNER JOIN transaction_tree tt ON st_map.tracked_transaction_id = tt.id
				WHERE (tt.beef_hex IS NULL OR st.beef_hex IS NOT NULL) 
			)
			SELECT DISTINCT * FROM transaction_tree
			ORDER BY id;
		`, uniqueIds).Scan(&transactions).Error

	if err != nil {
		return nil, err
	}

	// Convert to TxQueryResultSlice
	var slice beef.TxQueryResultSlice
	for _, record := range transactions {
		slice = append(slice, record.ToTxQueryResult())
	}
	return slice, nil
}

type stringSlice []string

func (ss stringSlice) unique() []string {
	set := make(map[string]struct{})
	for _, s := range ss {
		set[s] = struct{}{}
	}
	keys := maps.Keys(set)
	return slices.AppendSeq(make([]string, 0, len(set)), keys)
}
