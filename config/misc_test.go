package config

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandomHex will test the method RandomHex()
func TestRandomHex(t *testing.T) {
	t.Parallel()

	t.Run("valid tests", func(t *testing.T) {
		var tests = []struct {
			input          int
			expectedLength int
		}{
			{0, 0},
			{1, 2},
			{100000, 200000},
			{16, 32},
			{32, 64},
			{8, 16},
		}
		for _, test := range tests {
			output, err := RandomHex(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedLength, len(output))
		}
	})

	t.Run("panic tests", func(t *testing.T) {
		var tests = []struct {
			input          int
			expectedLength int
		}{
			{math.MaxInt, 16},
			{math.MaxInt - 1, 16},
		}
		for _, test := range tests {
			assert.Panics(t, func() {
				output, err := RandomHex(test.input)
				assert.Error(t, err)
				assert.Equal(t, test.expectedLength, len(output))
			})
		}
	})
}
