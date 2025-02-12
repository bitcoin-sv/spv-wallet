package repository

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
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

// HasTransactionInputSources checks if any of the given transaction inputs have a source transaction in the database.
// It queries the database to determine if at least one of the provided input source transaction IDs exists.
func (t *Transactions) HasTransactionInputSources(ctx context.Context, inputs ...*trx.TransactionInput) (bool, error) {
	sourceTXIDs := make([]string, 0, len(inputs))
	for _, input := range inputs {
		sourceTXIDs = append(sourceTXIDs, input.SourceTXID.String())
	}

	var count int64
	err := t.db.
		Model(&database.TrackedTransaction{}).
		WithContext(ctx).
		Where("id IN (?)", sourceTXIDs).
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, spverrors.Wrapf(err, "database query failed for source transactions %v", sourceTXIDs)
	}

	return count > 0, nil
}

// QueryTransactionInputSources retrieves the full ancestry of input sources for a given transaction.
// It recursively traces input sources in batches to minimize database queries.
func (t *Transactions) QueryTransactionInputSources(ctx context.Context, tx *trx.Transaction) (beef.TxQueryResultSlice, error) {
	visited := make(visitedTransactions)
	txID := tx.TxID().String()

	sourceTXID := make([]string, 0, len(tx.Inputs))
	for _, input := range tx.Inputs {
		sourceTXID = append(sourceTXID, input.SourceTXID.String())
	}

	// Fetch input sources in batches
	total, err := t.queryInputSourcesBatch(ctx, sourceTXID, visited)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to query input sources for transaction %s", txID)
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
		if !visitedTransactions.isVisited(txID) {
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
	results := make([]database.TrackedTransaction, 0, len(rows))
	var nextBatch []string
	for _, row := range rows {
		results = append(results, row)
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

	return append(results, nextResults...), nil
}

// visitedTransactions is a helper type that tracks visited transactions during recursive queries.
type visitedTransactions map[string]bool

// isVisited checks whether a transaction ID has already been visited.
// Returns true if the transaction ID has been processed, false otherwise.
func (v visitedTransactions) isVisited(txID string) bool {
	_, ok := v[txID]
	return ok
}

// recordVisited marks a transaction ID as visited to prevent redundant queries.
func (v visitedTransactions) recordVisited(txID string) {
	v[txID] = true
}
