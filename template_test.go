package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGreet will test the method Greet()
func TestGreet(t *testing.T) {
	got := Greet()

	assert.Equal(t, "Hi!", got, "should properly greet")
}
