package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/inputs"
	"gorm.io/gorm"
)

type InputsSelectorFixture interface {
	testabilities.DatabaseFixture
	NewInputSelector() inputs.Selector
}

type inputsSelectorFixture struct {
	testabilities.DatabaseFixture
	db *gorm.DB
}

func newFixture(t testing.TB) (InputsSelectorFixture, func()) {
	givenDB, cleanup := testabilities.Given(t)
	return &inputsSelectorFixture{
		DatabaseFixture: givenDB,
		db:              givenDB.GormDB(),
	}, cleanup
}

func (i *inputsSelectorFixture) NewInputSelector() inputs.Selector {
	return inputs.NewSelector(i.db, fixtures.DefaultFeeUnit)
}
