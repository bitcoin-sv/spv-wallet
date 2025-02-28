package txtestability

import "testing"

func New(t testing.TB) (given TransactionsFixtures) {
	given = Given(t)
	return
}
