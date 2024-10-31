package record_test

import (
	"context"
	"iter"
	"maps"
	"slices"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type mockRepository struct {
	transactions map[string]database.Transaction
	outputs      map[string]database.Output
	data         map[string]database.Data
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		transactions: make(map[string]database.Transaction),
		outputs:      make(map[string]database.Output),
		data:         make(map[string]database.Data),
	}
}

func (m *mockRepository) SaveTX(_ context.Context, txTable *database.Transaction, outputs []*database.Output, data []*database.Data) error {
	m.transactions[txTable.ID] = *txTable
	for _, output := range outputs {
		m.outputs[output.Outpoint().String()] = *output
	}
	for _, d := range data {
		m.data[d.Outpoint().String()] = *d
	}
	return nil
}

func (m *mockRepository) GetOutputs(_ context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error) {
	var outputs []*database.Output
	for outpoint := range outpoints {
		if output, ok := m.outputs[outpoint.String()]; ok {
			outputs = append(outputs, &output)
		}
	}
	return outputs, nil
}

func (m *mockRepository) withOutput(output database.Output) *mockRepository {
	m.outputs[output.Outpoint().String()] = output
	return m
}

func (m *mockRepository) getAllOutputs() []database.Output {
	return slices.Collect(maps.Values(m.outputs))
}

func (m *mockRepository) getAllData() []database.Data {
	return slices.Collect(maps.Values(m.data))
}

type mockBroadcaster struct {
	broadcastedTxs map[string]*trx.Transaction
}

func (m *mockBroadcaster) Broadcast(_ context.Context, tx *trx.Transaction) error {
	m.broadcastedTxs[tx.TxID().String()] = tx
	return nil
}

func newMockBroadcaster() *mockBroadcaster {
	return &mockBroadcaster{
		broadcastedTxs: make(map[string]*trx.Transaction),
	}
}
