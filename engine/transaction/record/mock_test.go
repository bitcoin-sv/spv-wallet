package record_test

import (
	"context"
	"fmt"
	"iter"

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

func (m *mockRepository) outpointID(txID string, vout uint32) string {
	return fmt.Sprintf("%s-%d", txID, vout)
}

func (m *mockRepository) SaveTX(_ context.Context, txTable *database.Transaction, outputs []*database.Output, data []*database.Data) error {
	m.transactions[txTable.ID] = *txTable
	for _, output := range outputs {
		m.outputs[m.outpointID(output.TxID, output.Vout)] = *output
	}
	for _, d := range data {
		m.data[m.outpointID(d.TxID, d.Vout)] = *d
	}
	return nil
}

func (m *mockRepository) GetOutputs(_ context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error) {
	var outputs []*database.Output
	for outpoint := range outpoints {
		key := m.outpointID(outpoint.TxID, outpoint.Vout)
		if output, ok := m.outputs[key]; ok {
			outputs = append(outputs, &output)
		}
	}
	return outputs, nil
}

func (m *mockRepository) withOutput(output database.Output) *mockRepository {
	m.outputs[m.outpointID(output.TxID, output.Vout)] = output
	return m
}

func (m *mockRepository) getOutput(txID string, vout uint32) (*database.Output, bool) {
	output, ok := m.outputs[m.outpointID(txID, vout)]
	return &output, ok
}

func (m *mockRepository) getData(txID string, vout uint32) (*database.Data, bool) {
	data, ok := m.data[m.outpointID(txID, vout)]
	return &data, ok
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
