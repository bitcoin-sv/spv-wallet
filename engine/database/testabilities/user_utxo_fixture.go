package testabilities

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"gorm.io/gorm"
)

var FirstCreatedAt = time.Date(2006, 02, 01, 15, 4, 5, 7, time.UTC)

type userUtxoFixture struct {
	db                 *gorm.DB
	t                  testing.TB
	index              uint
	userID             string
	txID               string
	vout               uint32
	satoshis           bsv.Satoshis
	estimatedInputSize uint64
}

func newUtxoFixture(t testing.TB, db *gorm.DB, index uint32) *userUtxoFixture {
	return &userUtxoFixture{
		t:                  t,
		db:                 db,
		index:              uint(index),
		userID:             fixtures.Sender.ID(),
		txID:               txIDTemplated(uint(index)),
		vout:               index,
		satoshis:           1,
		estimatedInputSize: database.EstimatedInputSizeForP2PKH,
	}
}

func txIDTemplated(index uint) string {
	return fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", index)
}

func (f *userUtxoFixture) OwnedBySender() UserUtxoFixture {
	f.userID = fixtures.Sender.ID()
	return f
}

func (f *userUtxoFixture) P2PKH() UserUtxoFixture {
	f.estimatedInputSize = database.EstimatedInputSizeForP2PKH
	return f
}

func (f *userUtxoFixture) WithSatoshis(satoshis bsv.Satoshis) UserUtxoFixture {
	f.satoshis = satoshis
	return f
}

func (f *userUtxoFixture) Stored() *database.UserUTXO {
	utxo := &database.UserUTXO{
		UserID:             f.userID,
		TxID:               f.txID,
		Vout:               f.vout,
		Satoshis:           uint64(f.satoshis),
		EstimatedInputSize: f.estimatedInputSize,
		Bucket:             string(bucket.BSV),
		CreatedAt:          FirstCreatedAt.Add(time.Duration(f.index) * time.Second), //nolint:gosec // this is used for testing and it should be fine even in case of integer overflow.
		TouchedAt:          FirstCreatedAt.Add(time.Duration(24) * time.Hour),
	}

	f.db.Create(utxo)

	return utxo
}
