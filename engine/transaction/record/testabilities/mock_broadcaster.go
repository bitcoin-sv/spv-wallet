package testabilities

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

type MockBroadcaster struct {
	broadcastedTxs map[string]*trx.Transaction
	returnErr      error
}

func NewMockBroadcaster() *MockBroadcaster {
	return &MockBroadcaster{
		broadcastedTxs: make(map[string]*trx.Transaction),
	}
}

func (m *MockBroadcaster) Broadcast(_ context.Context, tx *trx.Transaction) error {
	m.broadcastedTxs[tx.TxID().String()] = tx
	return m.returnErr
}

func (m *MockBroadcaster) willFailOnBroadcast(err error) *MockBroadcaster {
	m.returnErr = err
	return m
}

func (m *MockBroadcaster) checkBroadcasted(txID string) *trx.Transaction {
	tx := m.broadcastedTxs[txID]
	return tx
}
