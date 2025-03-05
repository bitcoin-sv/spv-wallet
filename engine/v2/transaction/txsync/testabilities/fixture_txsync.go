package testabilities

import (
	"testing"
	"time"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txsync"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const mockBlockHash = "00000000000000000f0905597b6cac80031f0f56834e74dce1a714c682a9ed38"
const mockBlockHeight = 885803

type FixtureTXsync interface {
	Service() *txsync.Service
	TXInfo(status chainmodels.TXStatus) TXInfoFixture
	MinedTXInfo() TXInfoFixture
	EmptyTXInfo() TXInfoFixture
	Repo() RepoFixtures
}

func Given(t testing.TB) FixtureTXsync {
	repo := newMockRepo(t)
	return &fixtureTXsync{
		t:       t,
		repo:    repo,
		service: txsync.NewService(tester.Logger(t), repo),
		givenTx: txtestability.Given(t),
	}
}

type RepoFixtures interface {
	ContainsBroadcastedTx(format Format) RepoFixtures
	WillFailOnGet() RepoFixtures
	WillFailOnUpdate() RepoFixtures
}

type fixtureTXsync struct {
	t       testing.TB
	givenTx txtestability.TransactionsFixtures
	service *txsync.Service
	repo    *MockRepo
}

func (f *fixtureTXsync) Service() *txsync.Service {
	return f.service
}

func (f *fixtureTXsync) TXInfo(status chainmodels.TXStatus) TXInfoFixture {
	require.NotNil(f.t, f.repo.subjectTx, "Test subject transaction is not set")

	return TXInfoFixture{
		TxID:      f.repo.subjectTx.ID(),
		Timestamp: time.Now().Add(10 * time.Minute),
		TXStatus:  status,
	}
}

func (f *fixtureTXsync) MinedTXInfo() TXInfoFixture {
	info := f.TXInfo(chainmodels.Mined)
	info.BlockHeight = mockBlockHeight
	info.MerklePath = mockBump(info.TxID).Hex()
	info.BlockHash = mockBlockHash

	return info
}

func (f *fixtureTXsync) EmptyTXInfo() TXInfoFixture {
	return TXInfoFixture{}
}

func (f *fixtureTXsync) Repo() RepoFixtures {
	return f.repo
}

type TXInfoFixture chainmodels.TXInfo

func (f TXInfoFixture) WithTimestamp(tim time.Time) TXInfoFixture {
	f.Timestamp = tim
	return f
}

func (f TXInfoFixture) WithWrongMerklePath() TXInfoFixture {
	f.MerklePath = "wrong"
	return f
}

func (f TXInfoFixture) WithSubjectTxOutOfMerklePath() TXInfoFixture {
	f.MerklePath = mockBump("84d44ab896962f9e57970af36cc4a17da3c8ce550f449e74641645074dc95622").Hex()
	return f
}

func mockBump(txID string) *trx.MerklePath {
	someHash, _ := chainhash.NewHashFromHex("4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06")
	txIDHash, _ := chainhash.NewHashFromHex(txID)
	bump := trx.NewMerklePath(mockBlockHeight, [][]*trx.PathElement{{
		&trx.PathElement{
			Hash:   someHash,
			Offset: 0,
		},
		&trx.PathElement{
			Hash:   txIDHash,
			Offset: 1,
			Txid:   lo.ToPtr(true),
		},
	}})

	return bump
}
