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
		pki, derivationKey, err := PaymailPKI(pubKey, "alice", "example.com")

		// then:
		assert.NoError(t, err)
		assert.Equal(t, "03a5399ee8ffff501739cb8167164ae88dcac6d1ca07cb863691dde12ae54012b8", pki.ToDERHex())
		assert.Equal(t, "1-paymail_pki-alice@example.com_0", derivationKey)
	})

	t.Run("try to generate PKI on nil", func(t *testing.T) {
		// when:
		pki, derivationKey, err := PaymailPKI(nil, "alice", "example.com")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
		assert.Equal(t, "1-paymail_pki-alice@example.com_0", derivationKey)
	})

	t.Run("try to generate PKI on empty alias", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		pki, derivationKey, err := PaymailPKI(pubKey, "", "example.com")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
		assert.Empty(t, derivationKey)
	})

	t.Run("try to generate PKI on empty domain", func(t *testing.T) {
		// given:
		pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

		// when:
		pki, derivationKey, err := PaymailPKI(pubKey, "alice", "")

		// then:
		assert.ErrorIs(t, err, ErrDeriveKey)
		assert.Nil(t, pki)
		assert.Empty(t, derivationKey)
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
