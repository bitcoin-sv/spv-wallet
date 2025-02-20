package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
)

func New(t testing.TB) (given TransactionsEndpointFixture, then testabilities.TransactionsEndpointAssertions) {
	given = Given(t)
	then = testabilities.NewTransactionsEndpointAssertions(t, given)
	return
}
