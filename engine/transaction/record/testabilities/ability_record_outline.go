package testabilities

import "testing"

func New(t testing.TB) (RecordServiceFixture, RecordOutlineAssert) {
	g := given(t)
	return g, then(t, g)
}
