package testabilities

import "testing"

func New(t testing.TB) (given FixtureTXsync, then AssertTXsync) {
	given = Given(t)
	then = Then(t, given)
	return
}
