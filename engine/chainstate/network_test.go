package chainstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNetwork_String will test the method String()
func TestNetwork_String(t *testing.T) {
	t.Parallel()

	t.Run("test all networks", func(t *testing.T) {
		assert.Equal(t, mainNet, MainNet.String())
		assert.Equal(t, stn, StressTestNet.String())
		assert.Equal(t, testNet, TestNet.String())
	})

	t.Run("unknown network", func(t *testing.T) {
		un := Network("")
		assert.Equal(t, "", un.String())
	})
}

// TestNetwork_Alternate will test the method Alternate()
func TestNetwork_Alternate(t *testing.T) {
	t.Parallel()

	t.Run("test all networks", func(t *testing.T) {
		assert.Equal(t, mainNetAlt, MainNet.Alternate())
		assert.Equal(t, stn, StressTestNet.Alternate())
		assert.Equal(t, testNetAlt, TestNet.Alternate())
	})

	t.Run("unknown network", func(t *testing.T) {
		un := Network("")
		assert.Equal(t, "", un.Alternate())
	})
}
