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
	set := make(map[string]struct{})
	for _, txID := range sourceTXIDs {
		set[txID] = struct{}{}
	}

	keys := maps.Keys(set)
	ids := slices.AppendSeq(make([]string, 0, len(set)), keys)

	var count int64
	err := t.db.
		Model(&database.TrackedTransaction{}).
		WithContext(ctx).
		Where("id IN (?)", ids).
		Count(&count).Error
	if err != nil {
		return false, spverrors.Wrapf(err, "database query failed for source transactions %v", sourceTXIDs)
	}

	return count == int64(len(ids)), nil
}

// FindTransactionInputSources retrieves the full ancestry of input sources for a given transaction.
// It recursively traces input sources in batches to minimize database queries.
func (t *Transactions) FindTransactionInputSources(ctx context.Context, sourceTXIDs ...string) (beef.TxQueryResultSlice, error) {
	visited := make(visitedTransactions)

	// Fetch input sources in batches
	total, err := t.queryInputSourcesBatch(ctx, sourceTXIDs, visited)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to query input sources batches")
	}

	// Convert to TxQueryResultSlice
	var slice beef.TxQueryResultSlice
	for _, record := range total {
		slice = append(slice, record.ToTxQueryResult())
	}
	return slice, nil
}

// queryInputSourcesBatch retrieves transactions in batches to optimize database performance.
// It avoids redundant queries by tracking visited transactions and continues retrieving input sources recursively.
func (t *Transactions) queryInputSourcesBatch(ctx context.Context, txIDs []string, visitedTransactions visitedTransactions) ([]database.TrackedTransaction, error) {
	if len(txIDs) == 0 {
		return nil, nil
	}

	// Filter out already visited transactions
	filteredIDs := make([]string, 0, len(txIDs))
	for _, txID := range txIDs {
		if visitedTransactions.isNotVisited(txID) {
			visitedTransactions.recordVisited(txID)
			filteredIDs = append(filteredIDs, txID)
		}
	}

	if len(filteredIDs) == 0 {
		return nil, nil
	}

	// Batch query transactions
	var rows []database.TrackedTransaction
	err := t.db.
		WithContext(ctx).
		Preload("SourceTxInputs").
		Where("id IN (?)", filteredIDs).
		Find(&rows).Error
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to batch query transactions: %v", filteredIDs)
	}

	// Process results and collect next batch of transaction IDs
	var nextBatch []string
	for _, row := range rows {
		if !row.HasBeefHex() {
			for _, input := range row.SourceTxInputs {
				nextBatch = append(nextBatch, input.SourceTxID)
			}
		}
	}

	// Recursively fetch the next batch
	nextResults, err := t.queryInputSourcesBatch(ctx, nextBatch, visitedTransactions)
	if err != nil {
		return nil, err
	}

	return append(rows, nextResults...), nil
}

type visitedTransactions map[string]struct{}

func (v visitedTransactions) isNotVisited(txID string) bool {
	_, ok := v[txID]
	return !ok
}

func (v visitedTransactions) recordVisited(txID string) { v[txID] = struct{}{} }
