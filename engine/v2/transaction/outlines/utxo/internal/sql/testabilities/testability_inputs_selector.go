package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
)

const SizeOfTransactionWithOnlyP2PKHOutput = 44

// MaxSizeWithoutFeeForSingleInput is the maximum size of a transaction that can be created without a fee for a single P2PKH input.
//
// We're calculating it by taking Fee Unit bytes and subtracting the estimated size of unlocking script for P2PKH
// because this script will be added to transaction after collecting inputs, and it will increase the size of the transaction,
// so it would impact the fee calculation (that's why we need to take it into account).
var MaxSizeWithoutFeeForSingleInput = fixtures.DefaultFeeUnit.Bytes - database.EstimatedInputSizeForP2PKH

func New(t testing.TB) (given InputsSelectorFixture, then InputsSelectorAssertions, cleanup func()) {
	given, cleanup = newFixture(t)
	then = newAssertions(t)
	return
}
