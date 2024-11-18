package testabilities

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/record"
)

type RecordServiceFixture interface {
	NewRecordService() *record.Service

	Repository() RepositoryFixture
	Broadcaster() BroadcasterFixture
}

type RepositoryFixture interface {
	WithOutputs(outputs ...database.Output) RepositoryFixture
	WithUTXOs(outpoints ...bsv.Outpoint) RepositoryFixture
	WillFailOnSaveTX(err error) RepositoryFixture
	WillFailOnGetOutputs(err error) RepositoryFixture
}

type BroadcasterFixture interface {
	WillFailOnBroadcast(err error) BroadcasterFixture
}

type recordServiceFixture struct {
	repository  *mockRepository
	broadcaster *mockBroadcaster
	t           testing.TB

	initialOutputs []database.Output
	initialData    []database.Data
}

func given(t testing.TB) *recordServiceFixture {
	return &recordServiceFixture{
		t:           t,
		repository:  newMockRepository(),
		broadcaster: newMockBroadcaster(),
	}
}

func (f *recordServiceFixture) Repository() RepositoryFixture {
	return f.repository
}

func (f *recordServiceFixture) Broadcaster() BroadcasterFixture {
	return f.broadcaster
}

func (f *recordServiceFixture) NewRecordService() *record.Service {
	f.initialOutputs = f.repository.GetAllOutputs()
	f.initialData = f.repository.GetAllData()

	return record.NewService(tester.Logger(f.t), f.repository, f.broadcaster)
}
