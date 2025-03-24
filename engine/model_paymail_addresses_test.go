package engine

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testXpubAuth = "xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA"
)

func TestNewPaymail(t *testing.T) {
	t.Run("paymail basic test", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, true, false)
		defer deferMe()

		paymail := "paymail@tester.com"
		externalDerivationNum := randDerivationNum()

		xPub, err := compat.GetHDKeyFromExtendedPublicKey(testXpubAuth)
		require.NoError(t, err)
		require.NotNil(t, xPub)

		// Get the external public key
		var paymailExternalKey *compat.ExtendedKey
		paymailExternalKey, err = compat.GetHDKeyByPath(
			xPub, utils.ChainExternal, externalDerivationNum,
		)
		require.NoError(t, err)
		require.NotNil(t, paymailExternalKey)

		var paymailIdentityKey *compat.ExtendedKey
		paymailIdentityKey, err = compat.GetHDKeyChild(paymailExternalKey, uint32(utils.MaxInt32))
		require.NoError(t, err)
		require.NotNil(t, paymailIdentityKey)

		paymailExternalXPub := paymailExternalKey.String()
		paymailIdentityXPub := paymailIdentityKey.String()

		p := newPaymail(
			paymail,
			externalDerivationNum,
			WithClient(client),
			WithXPub(testXpubAuth),
			WithEncryptionKey(testEncryption),
		)
		p.PublicName = "Tester"
		p.Avatar = "img url"
		err = p.Save(ctx)
		require.NoError(t, err)

		p2 := newPaymail(
			paymail,
			externalDerivationNum,
			WithClient(client),
			WithEncryptionKey(testEncryption),
		)
		p2.ID = "" // Remove ID (to make query work)
		conditions := map[string]interface{}{
			aliasField:  p.Alias,
			domainField: p.Domain,
		}
		err = Get(ctx, p2, conditions, false, 0, false)
		require.NoError(t, err)

		var identityKey *compat.ExtendedKey
		identityKey, err = p2.GetIdentityXpub()
		require.NoError(t, err)
		require.NotNil(t, identityKey)

		assert.Equal(t, paymail, p2.Alias+"@"+p2.Domain)
		assert.Equal(t, "Tester", p2.PublicName)
		assert.Equal(t, "img url", p2.Avatar)
		assert.Equal(t, "d8c2bed524071d72d859caf90da5f448b5861cd4d4fd47697f94166c13c5a987", p2.XpubID)
		assert.Equal(t, paymailIdentityXPub, identityKey.String())

		// Decrypt the external key
		var decrypted string
		decrypted, err = utils.Decrypt(testEncryption, p2.ExternalXpubKey)
		require.NoError(t, err)
		assert.Equal(t, paymailExternalXPub, decrypted)
	})

	t.Run("test derive child keys", func(t *testing.T) {
		// this is used in paymail to store the derived External xpub only in the DB
		hdKey, err := utils.ValidateXPub(testXPub)
		require.NoError(t, err)

		var internal, external string
		external, internal, err = utils.DeriveAddresses(
			hdKey, 1,
		)
		require.NoError(t, err)

		addressExternal := "16fq7PmmXXbFUG5maT5Xvr2zDBUgN1xdMF"
		addressInternal := "1PQW54xMn5KA6uK7wgfzN4y7ZXMi6o7Qtm"
		assert.Equal(t, addressInternal, internal)
		assert.Equal(t, addressExternal, external)

		childKeyChain0, _ := compat.GetHDKeyChild(hdKey, 0)
		childKeyChain01, _ := compat.GetHDKeyChild(childKeyChain0, 1)
		key0, _ := childKeyChain01.ECPubKey()
		address0, _ := script.NewAddressFromPublicKey(key0, true)
		assert.Equal(t, addressExternal, address0.AddressString)

		childKeyChain1, _ := compat.GetHDKeyChild(hdKey, 1)
		childKeyChain11, _ := compat.GetHDKeyChild(childKeyChain1, 1)
		key1, _ := childKeyChain11.ECPubKey()
		address1, _ := script.NewAddressFromPublicKey(key1, true)
		assert.Equal(t, addressInternal, address1.AddressString)
	})

	t.Run("ctor test", func(t *testing.T) {
		// given
		derivationNumber := randDerivationNum()

		// when
		pm := newPaymail("paymail@domain.sc", derivationNumber, WithXPub(testXPub))

		// then
		require.Equal(t, uint32(0), pm.PubKeyNum, "PubKeyNum for new paymail MUST equal 0")
		require.Equal(t, uint32(0), pm.XpubDerivationSeq, "XpubDerivationSeq for new paymail MUST equal 0")
		require.Equal(t, derivationNumber, pm.ExternalXpubKeyNum)
	})

	t.Run("GetNextXpub() test", func(t *testing.T) {
		ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithFreeCache())
		defer deferMe()

		// given
		derivationNumber := randDerivationNum()

		pm := newPaymail("paymail@domain.sc", derivationNumber, WithClient(c), WithXPub(testXPub))
		err := pm.Save(ctx)
		require.NoError(t, err)

		// when
		firstExternalXpub, err := pm.GetNextXpub(ctx)
		require.NoError(t, err)
		derivationForFirstXpub := pm.XpubDerivationSeq

		secondExternalXpub, err := pm.GetNextXpub(ctx)
		require.NoError(t, err)
		derivationForSecondXpub := pm.XpubDerivationSeq

		// then
		require.Equal(t, derivationNumber, pm.ExternalXpubKeyNum, "ExternalXpubKeyNum MUST NOT be changed during any key rotation")

		require.Equal(t, uint32(1), derivationForFirstXpub, "XpubDerivationSeq after first rotation MUST equal 1")
		require.Equal(t, uint32(2), derivationForSecondXpub, "XpubDerivationSeq after second rotation MUST equal 2")

		require.NotEqual(t, firstExternalXpub, secondExternalXpub, "External Xpubs cannot be equal")

		// verify the correctness of the derivation path
		expectedXpubDerivationPath := fmt.Sprintf("%d/%d/%d", utils.ChainExternal, derivationNumber, pm.XpubDerivationSeq)
		masterXPub, _ := compat.GetHDKeyFromExtendedPublicKey(testXPub)

		expectedExternalXpub, err := masterXPub.DeriveChildFromPath(expectedXpubDerivationPath)
		require.NoError(t, err)

		require.Equal(t, expectedExternalXpub, secondExternalXpub)
	})

	t.Run("GetPubKey() test", func(t *testing.T) {
		ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithFreeCache())
		defer deferMe()

		// given
		derivationNumber := randDerivationNum()

		pm := newPaymail("paymail@domain.sc", derivationNumber, WithClient(c), WithXPub(testXPub))
		err := pm.Save(ctx)
		require.NoError(t, err)

		initialDerivationSeq := pm.XpubDerivationSeq

		// when
		firstPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		secondPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		// then
		require.Equal(t, initialDerivationSeq, pm.XpubDerivationSeq, "XpubDerivationSeq cannot be changed")
		require.Equal(t, derivationNumber, pm.ExternalXpubKeyNum, "ExternalXpubKeyNum MUST NOT be changed during any key rotation")

		require.Equal(t, firstPubKey, secondPubKey, "PubKeys must be equal")

		// verify the correctness of the derivation path
		expectedXpubDerivationPath := fmt.Sprintf("%d/%d/%d", utils.ChainExternal, derivationNumber, initialDerivationSeq)
		masterXPub, _ := compat.GetHDKeyFromExtendedPublicKey(testXPub)

		expectedHdPubKey, err := masterXPub.DeriveChildFromPath(expectedXpubDerivationPath)
		require.NoError(t, err)

		expectedPubKey, err := expectedHdPubKey.ECPubKey()
		require.NoError(t, err)

		require.Equal(t, hex.EncodeToString(expectedPubKey.Compressed()), firstPubKey)
	})

	t.Run("RotatePubKey() test", func(t *testing.T) {
		ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithFreeCache())
		defer deferMe()

		// given
		derivationNumber := randDerivationNum()

		pm := newPaymail("paymail@domain.sc", derivationNumber, WithClient(c), WithXPub(testXPub))
		err := pm.Save(ctx)
		require.NoError(t, err)

		initialDerivationSeq := pm.XpubDerivationSeq

		// when
		firstPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		err = pm.RotatePubKey(ctx)
		require.NoError(t, err)

		secondPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		// then
		require.Greater(t, pm.XpubDerivationSeq, initialDerivationSeq, "XpubDerivationSeq must be incremented after rotation")
		require.Equal(t, derivationNumber, pm.ExternalXpubKeyNum, "ExternalXpubKeyNum MUST NOT be changed during any key rotation")

		require.NotEqual(t, firstPubKey, secondPubKey)

		externalXPubDerivationPath := fmt.Sprintf("%d/%d/%d", utils.ChainExternal, derivationNumber, pm.XpubDerivationSeq)
		masterXPub, _ := compat.GetHDKeyFromExtendedPublicKey(testXPub)

		expectedExternalXpub, err := masterXPub.DeriveChildFromPath(externalXPubDerivationPath)
		require.NoError(t, err)

		expectedPubKey, err := expectedExternalXpub.ECPubKey()
		require.NoError(t, err)

		require.Equal(t, hex.EncodeToString(expectedPubKey.Compressed()), secondPubKey)
	})

	t.Run("ExternalXPub and PubKey rotation test", func(t *testing.T) {
		ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithFreeCache())
		defer deferMe()

		// given
		derivationNumber := randDerivationNum()

		pm := newPaymail("paymail@domain.sc", derivationNumber, WithClient(c), WithXPub(testXPub))
		err := pm.Save(ctx)
		require.NoError(t, err)

		initialDerivationSeq := pm.XpubDerivationSeq

		// when

		// get pub key	- XpubDerivationSeq - should not be changed
		// PubKeyNum shoud be equal to XpubDerivationSeq
		firstPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		require.Equal(t, initialDerivationSeq, pm.XpubDerivationSeq)
		require.Equal(t, pm.PubKeyNum, pm.XpubDerivationSeq)

		// get external xpub - XpubDerivationSeq should be incremented
		firstExternalXpub, err := pm.GetNextXpub(ctx)
		require.NoError(t, err)

		require.Greater(t, pm.XpubDerivationSeq, initialDerivationSeq)

		// get pub key again - should be the same as previous one
		secondPubKey, err := pm.GetPubKey()
		require.NoError(t, err)

		require.Equal(t, firstPubKey, secondPubKey)

		// rotate pub key - XpubDerivationSeq should be incremented
		err = pm.RotatePubKey(ctx)
		require.NoError(t, err)

		// get external xpub - XpubDerivationSeq should be incremented
		secondExternalXpub, err := pm.GetNextXpub(ctx)
		require.NoError(t, err)

		numberOfRotation := 3
		require.Equal(t, initialDerivationSeq+uint32(numberOfRotation), pm.XpubDerivationSeq)
		require.NotEqual(t, firstExternalXpub, secondExternalXpub)
	})
}

func randDerivationNum() uint32 {
	rnd := rand.Int63n(int64(compat.HardenedKeyStart))
	if rnd < 0 {
		rnd = rnd * -1
	}

	return uint32(rnd)
}
