package testabilities

import "testing"

func New(t testing.TB) (ContactFixture, ContactAssertion) {
	return given(t), then(t)
}
