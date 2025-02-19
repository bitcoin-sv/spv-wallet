package type42

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestination(t *testing.T) {
	t.Run("try to generate random destination", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		dst, err := NewDestinationWithRandomReference(pubKey)

		// then:
		assert.NoError(t, err)
		assert.NotEmpty(t, dst.PubKey.ToDERHex())
		assert.Contains(t, dst.DerivationKey, "1-destination-")
		assert.NotEmpty(t, dst.ReferenceID)
	})
}
