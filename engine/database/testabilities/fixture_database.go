package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
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
	// P2PKH ensures that the UTXO is with P2PKH locking script (which is default behavior of this fixture).
	P2PKH() UserUtxoFixture
	// WithSatoshis sets the satoshis value of the UTXO.
	WithSatoshis(satoshis bsv.Satoshis) UserUtxoFixture

	Storable[database.UserUtxos]
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

func Given(t testing.TB) (given DatabaseFixture, cleanup func()) {
	engineWithConfign, cleanup := testengine.Given(t).Engine()

	db := engineWithConfign.Engine.Datastore().DB()
	fixture := &databaseFixture{
		t:  t,
		db: db,
	}

	// TODO: remove this when we will include UserUtxos in the production code
	err := fixture.db.AutoMigrate(&database.UserUtxos{})
	require.NoError(t, err)

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
