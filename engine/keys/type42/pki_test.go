package type42

import (
	"testing"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
)

func TestPKI(t *testing.T) {
	t.Run("generate PKI", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		pki, err := PKI(pubKey)

		// then:
		assert.NoError(t, err)
		assert.Equal(t, "02b9a822f2db22649e14eedf75ba140cf5dacc6b2690cfae9da55b551069461705", pki.ToDERHex())
	})

	t.Run("try to generate PKI on nil", func(t *testing.T) {
		// when:
		pki, err := PKI(nil)

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
	})
}

func makePubKey(t *testing.T, pubDERHex string) *primitives.PublicKey {
	t.Helper()
	pk, err := primitives.PublicKeyFromString(pubDERHex)
	if err != nil {
		t.Fatalf("failed to create public key: %s", err)
	}
	return pk
}
