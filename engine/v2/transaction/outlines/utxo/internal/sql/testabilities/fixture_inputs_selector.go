package testabilities

import (
	"slices"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo/internal/sql"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type InputsSelectorFixture interface {
	testabilities.DatabaseFixture
	NewInputSelector() *sql.UTXOSelector
	Transaction() InputsSelectorTransactionFixture
}

type InputsSelectorTransactionFixture interface {
	ForSatoshisAndSize(SatoshisAndSizeProvider) *sdk.Transaction
}

type SatoshisAndSizeProvider interface {
	Satoshis() bsv.Satoshis
	Size() int
}

type inputsSelectorFixture struct {
	testabilities.DatabaseFixture
	t           testing.TB
	db          *gorm.DB
	transaction txtestability.TransactionsFixtures
}

func newFixture(t testing.TB) (InputsSelectorFixture, func()) {
	givenDB, cleanup := testabilities.Given(t, testengine.WithV2())
	givenTx := txtestability.Given(t)
	return &inputsSelectorFixture{
		t:               t,
		DatabaseFixture: givenDB,
		db:              givenDB.GormDB(),
		transaction:     givenTx,
	}, cleanup
}

func (i *inputsSelectorFixture) NewInputSelector() *sql.UTXOSelector {
	return sql.NewUTXOSelector(i.db, fixtures.DefaultFeeUnit)
}

func (i *inputsSelectorFixture) Transaction() InputsSelectorTransactionFixture {
	return i
}

func (i *inputsSelectorFixture) ForSatoshisAndSize(provider SatoshisAndSizeProvider) *sdk.Transaction {
	sats := uint64(provider.Satoshis())
	size := provider.Size()

	bsvTransaction := i.transaction.Tx().
		WithP2PKHOutput(sats).
		TX()

	if size < bsvTransaction.Size() {
		i.t.Fatalf("Not implemented fixture for transaction with size smaller then for tx with only P2PKH output")
	}

	if size > bsvTransaction.Size() {
		addition := size - bsvTransaction.Size() - 2
		if addition < 0 {
			i.t.Fatalf("Not implemented fixture for transaction with size smaller then for tx with only P2PKH output + 1 bytes")
		}

		lockingScript := bsvTransaction.OutputIdx(0).LockingScript.Bytes()
		lockingScript = append(lockingScript, slices.Repeat([]byte{script.Op1}, addition)...)
		bsvTransaction.OutputIdx(0).LockingScript = script.NewFromBytes(lockingScript)
	}

	// Ensure that the fixture was created correctly
	require.Equal(i.t, size, bsvTransaction.Size(), "Failed to create transaction fixture with given size")
	require.Equal(i.t, sats, bsvTransaction.TotalOutputSatoshis(), "Failed to create transaction fixture with given satoshis")

	return bsvTransaction
}
