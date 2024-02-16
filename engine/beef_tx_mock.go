package engine

import (
	"context"
)

type MockTransactionStore struct {
	Transactions map[string]*Transaction
}

func NewMockTransactionStore() *MockTransactionStore {
	return &MockTransactionStore{
		Transactions: make(map[string]*Transaction),
	}
}

func (m *MockTransactionStore) AddToStore(tx *Transaction) {
	m.Transactions[tx.ID] = tx
}

func (m *MockTransactionStore) GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error) {
	var txs []*Transaction
	for _, txID := range txIDs {
		if tx, exists := m.Transactions[txID]; exists {
			txs = append(txs, tx)
		}
	}
	return txs, nil
}
