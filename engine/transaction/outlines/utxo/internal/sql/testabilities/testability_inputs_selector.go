package testabilities

import "testing"

func New(t testing.TB) (given InputsSelectorFixture, then InputsSelectorAssertions, cleanup func()) {
	given, cleanup = newFixture(t)
	then = newAssertions(t)
	return
}
