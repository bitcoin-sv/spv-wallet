package testabilities

import (
	"testing"
)

func New(t testing.TB) (given TransactionsEndpointFixture, then TransactionsEndpointAssertions) {
	given = Given(t)
	then = Then(t, given)
	return
}
