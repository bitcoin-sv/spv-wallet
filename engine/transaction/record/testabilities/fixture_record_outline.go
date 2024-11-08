package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type RecordServiceFixture interface {
	NewRecordService() *record.Service

	WithStoredUTXO(outpoints ...bsv.Outpoint) RecordServiceFixture
	WithStoredOutputs(outputs ...database.Output) RecordServiceFixture

	WillFailOnBroadcast(err error) RecordServiceFixture
}

type recordServiceFixture struct {
	repository  *MockRepository
	broadcaster *MockBroadcaster
	t           testing.TB

	initialOutputs []database.Output
	initialData    []database.Data
}

func given(t testing.TB) *recordServiceFixture {
	return &recordServiceFixture{
		t:           t,
		repository:  NewMockRepository(),
		broadcaster: NewMockBroadcaster(),
	}
}

func (f *recordServiceFixture) NewRecordService() *record.Service {
	f.initialOutputs = f.repository.GetAllOutputs()
	f.initialData = f.repository.GetAllData()

	return record.NewService(tester.Logger(f.t), f.repository, f.broadcaster)
}

func (f *recordServiceFixture) WithStoredUTXO(outpoints ...bsv.Outpoint) RecordServiceFixture {
	for _, outpoint := range outpoints {
		f.repository.WithUTXO(outpoint)
	}
	return f
}

func (f *recordServiceFixture) WithStoredOutputs(outputs ...database.Output) RecordServiceFixture {
	for _, output := range outputs {
		f.repository.WithOutput(output)
	}
	return f
}

func (f *recordServiceFixture) WillFailOnBroadcast(err error) RecordServiceFixture {
	f.broadcaster.willFailOnBroadcast(err)
	return f
}
