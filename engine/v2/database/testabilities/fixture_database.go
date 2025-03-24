package testabilities

import (
	"testing"

	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

type DatabaseFixture interface {
	DB() DatabaseDataFixtures
	GormDB() *gorm.DB
}

type DatabaseDataFixtures interface {
	HasUTXO() UserUtxoFixture
}

type UserUtxoFixture interface {
	// OwnedBySender ensures that the UTXO is owned by the sender (which is default behavior of this fixture).
	OwnedBySender() UserUtxoFixture
	// OwnedByRecipient ensures that the UTXO is owned by the recipient
	OwnedByRecipient() UserUtxoFixture
	// P2PKH ensures that the UTXO is with P2PKH locking script (which is default behavior of this fixture).
	P2PKH() UserUtxoFixture
	// WithSatoshis sets the satoshis value of the UTXO.
	WithSatoshis(satoshis bsv.Satoshis) UserUtxoFixture

	Storable[database.UserUTXO]
}

type Storable[Data any] interface {
	// Stored ensures that the data is stored in the database.
	Stored() *Data
}

type databaseFixture struct {
	t                testing.TB
	db               *gorm.DB
	utxoEntriesIndex uint32
}

func Given(t testing.TB, opts ...testengine.ConfigOpts) (given DatabaseFixture, cleanup func()) {
	engineWithConfig, cleanup := testengine.Given(t).EngineWithConfiguration(opts...)

	db := engineWithConfig.Engine.DB()
	fixture := &databaseFixture{
		t:  t,
		db: db,
	}

	return fixture, cleanup
}

func (f *databaseFixture) DB() DatabaseDataFixtures {
	return f
}

func (f *databaseFixture) GormDB() *gorm.DB {
	return f.db
}

func (f *databaseFixture) HasUTXO() UserUtxoFixture {
	f.utxoEntriesIndex++
	return newUtxoFixture(f.t, f.db, f.utxoEntriesIndex)
}
