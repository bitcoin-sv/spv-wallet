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
)

const mockBlockHash = "00000000000000000f0905597b6cac80031f0f56834e74dce1a714c682a9ed38"
const mockBlockHeight = 885803

type FixtureTXsync interface {
	Service() *txsync.Service
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
	WillFailOn(f FailingPoint)
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

func (f *fixtureTXsync) Repo() RepoFixtures {
	return f.repo
}

func TXInfo(t testing.TB, status chainmodels.TXStatus) TXInfoSpec {

	return TXInfoSpec{
		TxID:      MockTx(t).ID(),
		Timestamp: time.Now().Add(10 * time.Minute),
		TXStatus:  status,
	}
}

func MinedTXInfo(t testing.TB) TXInfoSpec {
	info := TXInfo(t, chainmodels.Mined)
	info.BlockHeight = mockBlockHeight
	info.MerklePath = mockBump(info.TxID).Hex()
	info.BlockHash = mockBlockHash

	return info
}

func EmptyTXInfo() TXInfoSpec {
	return TXInfoSpec{}
}

type TXInfoSpec chainmodels.TXInfo

func (f TXInfoSpec) WithTimestamp(tim time.Time) TXInfoSpec {
	f.Timestamp = tim
	return f
}

func (f TXInfoSpec) WithWrongMerklePath() TXInfoSpec {
	f.MerklePath = "wrong"
	return f
}

func (f TXInfoSpec) WithSubjectTxOutOfMerklePath() TXInfoSpec {
	f.MerklePath = mockBump("84d44ab896962f9e57970af36cc4a17da3c8ce550f449e74641645074dc95622").Hex()
	return f
}

func (f TXInfoSpec) IncrementBlockHeight() TXInfoSpec {
	f.BlockHeight++
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
