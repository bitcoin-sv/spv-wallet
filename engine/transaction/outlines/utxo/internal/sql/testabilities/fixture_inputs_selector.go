package testabilities

import (
	"testing"

	"gorm.io/gorm"

	"github.com/bitcoin-sv/spv-wallet/engine/database/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/utxo/internal/sql"
)

type InputsSelectorFixture interface {
	testabilities.DatabaseFixture
	NewInputSelector() sql.UTXOSelector
}

type inputsSelectorFixture struct {
	testabilities.DatabaseFixture
	db *gorm.DB
}

func newFixture(t testing.TB) (InputsSelectorFixture, func()) {
	givenDB, cleanup := testabilities.Given(t, testengine.WithV2())
	return &inputsSelectorFixture{
		DatabaseFixture: givenDB,
		db:              givenDB.GormDB(),
	}, cleanup
}

func (i *inputsSelectorFixture) NewInputSelector() sql.UTXOSelector {
	return sql.NewUTXOSelector(i.db, fixtures.DefaultFeeUnit)
}
