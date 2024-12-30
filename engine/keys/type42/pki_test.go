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
		pki, err := PaymailPKI(pubKey, "alice", "example.com")

		// then:
		assert.NoError(t, err)
		assert.Equal(t, "03a5399ee8ffff501739cb8167164ae88dcac6d1ca07cb863691dde12ae54012b8", pki.ToDERHex())
	})

	t.Run("try to generate PKI on nil", func(t *testing.T) {
		// when:
		pki, err := PaymailPKI(nil, "alice", "example.com")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
	})

	t.Run("try to generate PKI on empty alias", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		pki, err := PaymailPKI(pubKey, "", "example.com")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
	})

	t.Run("try to generate PKI on empty domain", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		pki, err := PaymailPKI(pubKey, "alice", "")

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
