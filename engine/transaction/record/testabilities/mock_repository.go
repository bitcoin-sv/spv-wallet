package testabilities

import (
	"context"
	"iter"
	"maps"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type MockRepository struct {
	transactions map[string]database.Transaction
	outputs      map[string]database.Output
	data         map[string]database.Data

	errOnSave       error
	errOnGetOutputs error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		transactions: make(map[string]database.Transaction),
		outputs:      make(map[string]database.Output),
		data:         make(map[string]database.Data),
	}
}

func (m *MockRepository) SaveTX(_ context.Context, txTable *database.Transaction, outputs []*database.Output, data []*database.Data) error {
	if m.errOnSave != nil {
		return m.errOnSave
	}
	m.transactions[txTable.ID] = *txTable
	for _, output := range outputs {
		m.outputs[output.Outpoint().String()] = *output
	}
	for _, d := range data {
		m.data[d.Outpoint().String()] = *d
	}
	return nil
}

func (m *MockRepository) GetOutputs(_ context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error) {
	if m.errOnGetOutputs != nil {
		return nil, m.errOnGetOutputs
	}
	var outputs []*database.Output
	for outpoint := range outpoints {
		if output, ok := m.outputs[outpoint.String()]; ok {
			outputs = append(outputs, &output)
		}
	}
	return outputs, nil
}

func (m *MockRepository) WithOutput(output database.Output) *MockRepository {
	m.outputs[output.Outpoint().String()] = output
	return m
}

func (m *MockRepository) WithUTXO(outpoint bsv.Outpoint) *MockRepository {
	m.outputs[outpoint.String()] = database.Output{
		TxID: outpoint.TxID,
		Vout: outpoint.Vout,
	}
	return m
}

func (m *MockRepository) WillFailOnSaveTX(err error) *MockRepository {
	m.errOnSave = err
	return m
}

func (m *MockRepository) WillFailOnGetOutputs(err error) *MockRepository {
	m.errOnGetOutputs = err
	return m
}

func (m *MockRepository) GetAllOutputs() []database.Output {
	return slices.Collect(maps.Values(m.outputs))
}

func (m *MockRepository) GetAllData() []database.Data {
	return slices.Collect(maps.Values(m.data))
}

func (m *MockRepository) getTransaction(txID string) *database.Transaction {
	tx, ok := m.transactions[txID]
	if !ok {
		return nil
	}
	return &tx
}
