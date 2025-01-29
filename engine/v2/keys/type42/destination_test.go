package type42

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestination(t *testing.T) {
	t.Run("generate destination", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")
		referenceID := "4c7bc22854691fda2d643f9c5cf6d218"

		// when:
		destination, derivationKey, err := Destination(pubKey, referenceID)

		// then:
		assert.NoError(t, err)
		assert.Equal(t, "03d34d33cb9cf83ad5bea49c7ebb1adafe0c85ceda0e256a0c5db3b6cb28e3ec99", destination.ToDERHex())
		assert.Equal(t, "1-destination-"+referenceID, derivationKey)
	})

	t.Run("try to generate destination on nil", func(t *testing.T) {
		// given:
		referenceID := "4c7bc22854691fda2d643f9c5cf6d218"

		// when:
		destination, derivationKey, err := Destination(nil, referenceID)

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, destination)
		assert.Equal(t, "1-destination-"+referenceID, derivationKey)
	})

	t.Run("try to generate destination on empty referenceID", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		destination, derivationKey, err := Destination(pubKey, "")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, destination)
		assert.Equal(t, "", derivationKey)
	})
}
