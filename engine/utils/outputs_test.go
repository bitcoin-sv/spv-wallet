package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitOutputValues(t *testing.T) {
	t.Parallel()

	t.Run("1 value", func(t *testing.T) {
		outputValues, err := SplitOutputValues(10000, 1)
		require.NoError(t, err)
		assert.Equal(t, []uint64{10000}, outputValues)
	})

	t.Run("2 values", func(t *testing.T) {
		outputValues, err := SplitOutputValues(10000, 2)
		require.NoError(t, err)
		assert.Len(t, outputValues, 2)
		totalOutput := uint64(0)
		for _, output := range outputValues {
			totalOutput += output
		}
		assert.Equal(t, uint64(10000), totalOutput)
	})

	t.Run("3 values", func(t *testing.T) {
		outputValues, err := SplitOutputValues(10000, 3)
		require.NoError(t, err)
		assert.Len(t, outputValues, 3)
		totalOutput := uint64(0)
		for _, output := range outputValues {
			totalOutput += output
		}
		assert.Equal(t, uint64(10000), totalOutput)
	})
}
