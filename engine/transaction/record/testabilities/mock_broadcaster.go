package testabilities

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

type mockBroadcaster struct {
	broadcastedTxs map[string]*trx.Transaction
	returnErr      error
}

func newMockBroadcaster() *mockBroadcaster {
	return &mockBroadcaster{
		broadcastedTxs: make(map[string]*trx.Transaction),
	}
}

func (m *mockBroadcaster) Broadcast(_ context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error) {
	m.broadcastedTxs[tx.TxID().String()] = tx
	return &chainmodels.TXInfo{
		TXStatus: chainmodels.SeenOnNetwork,
		TxID:     tx.TxID().String(),
	}, m.returnErr
}

func (m *mockBroadcaster) WillFailOnBroadcast(err error) BroadcasterFixture {
	m.returnErr = err
	return m
}

func (m *mockBroadcaster) checkBroadcasted(txID string) *trx.Transaction {
	tx := m.broadcastedTxs[txID]
	return tx
}
