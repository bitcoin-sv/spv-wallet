package testabilities

import (
	"testing"

	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo/internal/sql"
	"gorm.io/gorm"
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
	givenDB, cleanup := testabilities.Given(t, testengine.WithNewTransactionFlowEnabled())
	return &inputsSelectorFixture{
		DatabaseFixture: givenDB,
		db:              givenDB.GormDB(),
	}, cleanup
}

func (i *inputsSelectorFixture) NewInputSelector() sql.UTXOSelector {
	return sql.NewUTXOSelector(i.db, fixtures.DefaultFeeUnit)
}
