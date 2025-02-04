package testabilities

import (
	"testing"
)

func New(t testing.TB) (given TransactionsEndpointFixture, then TransactionsEndpointAssertions) {
	return Given(t), Then(t)
}
