package engine

import (
	"context"
)

// MockTransactionStore is a mock implementation of TransactionGetter.
type MockTransactionStore struct {
	Transactions map[string]*Transaction
}

// NewMockTransactionStore creates new mock transaction store.
func NewMockTransactionStore() *MockTransactionStore {
	return &MockTransactionStore{
		Transactions: make(map[string]*Transaction),
	}
}

// AddToStore adds transaction to store.
func (m *MockTransactionStore) AddToStore(tx *Transaction) {
	m.Transactions[tx.ID] = tx
}

// GetTransactionsByIDs returns transactions by IDs.
func (m *MockTransactionStore) GetTransactionsByIDs(_ context.Context, txIDs []string) ([]*Transaction, error) {
	var txs []*Transaction
	for _, txID := range txIDs {
		if tx, exists := m.Transactions[txID]; exists {
			txs = append(txs, tx)
		}
	}
	return txs, nil
}
