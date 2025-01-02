package testabilities

import (
	"context"
	"iter"
	"maps"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type mockRepository struct {
	transactions map[string]database.TrackedTransaction
	outputs      map[string]database.Output
	data         map[string]database.Data

	errOnSave       error
	errOnGetOutputs error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		transactions: make(map[string]database.TrackedTransaction),
		outputs:      make(map[string]database.Output),
		data:         make(map[string]database.Data),
	}
}

func (m *mockRepository) SaveTX(_ context.Context, txTable *database.TrackedTransaction) error {
	if m.errOnSave != nil {
		return m.errOnSave
	}
	m.transactions[txTable.ID] = *txTable
	for _, output := range txTable.Outputs {
		m.outputs[output.Outpoint().String()] = *output
	}
	for _, output := range txTable.Inputs {
		utxo := *output
		utxo.SpendingTX = txTable.ID
		m.outputs[utxo.Outpoint().String()] = utxo
	}
	for _, d := range txTable.Data {
		m.data[d.Outpoint().String()] = *d
	}
	return nil
}

func (m *mockRepository) GetOutputs(_ context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error) {
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

func (m *mockRepository) WithOutputs(outputs ...database.Output) RepositoryFixture {
	for _, output := range outputs {
		m.outputs[output.Outpoint().String()] = output
	}
	return m
}

func (m *mockRepository) WithUTXOs(outpoints ...bsv.Outpoint) RepositoryFixture {
	for _, o := range outpoints {
		m.outputs[o.String()] = database.Output{
			TxID: o.TxID,
			Vout: o.Vout,
		}
	}
	return m
}

func (m *mockRepository) WillFailOnSaveTX(err error) RepositoryFixture {
	m.errOnSave = err
	return m
}

func (m *mockRepository) WillFailOnGetOutputs(err error) RepositoryFixture {
	m.errOnGetOutputs = err
	return m
}

func (m *mockRepository) GetAllOutputs() []database.Output {
	return slices.Collect(maps.Values(m.outputs))
}

func (m *mockRepository) GetAllData() []database.Data {
	return slices.Collect(maps.Values(m.data))
}

func (m *mockRepository) getTransaction(txID string) *database.TrackedTransaction {
	tx, ok := m.transactions[txID]
	if !ok {
		return nil
	}
	return &tx
}

func (m *mockRepository) CheckAddress(_ context.Context, _ string) (bool, error) {
	return true, nil
}
