package utils

import (
	"encoding/hex"
	"strings"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExternalAddress = "1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt"
	testHash            = "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414"
	testInternalAddress = "14n4rKed7f5vkPfV7Yj8N3E8Pxa35Rytp9"
	testValue           = "1234567"
	testXpriv           = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	testXpub            = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	testXpubHash        = "1a0b10d4eda0636aae1709e7e7080485a4d99af3ca2962c6e677cf5b53d8ab8c"
	derivedXpriv        = "xprvA8mj2ZL1w6Nqpi6D2amJLo4Gxy24tW9uv82nQKmamT2rkg5DgjzJZRFnW33e7QJwn65uUWSuN6YQyWrujNjZdVShPRnpNUSRVTru4cxaqfd"
	derivedXpub         = "xpub6Mm5S4rumTw93CAg8cJJhw11WzrZHxsmHLxPCiBCKnZqdUQNEHJZ7DaGMKucRzXPHtoS2ZqsVSRjxVbibEvwmR2wXkZDd8RrTftmm42cRsf"
	privateKey0         = "5e319f45f94450c97aab649b263cfbad81485d265206548dab4c3046e26fcd03"
	privateKeyHash      = "644fdfc7e2815555e68d0317535b08d28a5e21d55d5c2d57e605cb63a346d9f2"
)

func TestHash(t *testing.T) {
	t.Parallel()

	t.Run("empty xpub", func(t *testing.T) {
		assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", Hash(""))
	})
	t.Run("example hash", func(t *testing.T) {
		assert.Equal(t, testHash, Hash(testValue))
	})
	t.Run("valid xpub test", func(t *testing.T) {
		assert.Equal(t, testXpubHash, Hash(testXpub))
	})
}

func TestRandomHex(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name           string
		input          int
		expectedLength int
	}{
		{"zero", 0, 0},
		{"one", 1, 2},
		{"100k", 100000, 200000},
		{"16->32", 16, 32},
		{"32->64", 32, 64},
		{"8->16", 8, 16},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := RandomHex(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedLength, len(output))
		})
	}
}

func TestValidateXPub(t *testing.T) {
	t.Parallel()

	t.Run("valid xpub", func(t *testing.T) {
		hdKey, err := ValidateXPub(testXpub)
		require.NotNil(t, hdKey)
		require.NoError(t, err)
	})

	t.Run("invalid length", func(t *testing.T) {
		hdKey, err := ValidateXPub("1234567890")
		assert.Nil(t, hdKey)
		assert.Error(t, err)
	})

	t.Run("unable to decode key", func(t *testing.T) {
		hdKey, err := ValidateXPub(strings.Replace(testXpub, "A", "B", -1))
		assert.Nil(t, hdKey)
		assert.Error(t, err)
	})
}

func TestDeriveAddresses(t *testing.T) {
	t.Parallel()

	t.Run("valid derive", func(t *testing.T) {
		hdKey, err := ValidateXPub(testXpub)
		require.NotNil(t, hdKey)
		require.NoError(t, err)

		var internal, external string
		external, internal, err = DeriveAddresses(hdKey, 0)
		require.NoError(t, err)
		assert.Equal(t, testExternalAddress, external)
		assert.Equal(t, testInternalAddress, internal)
	})

	t.Run("nil key", func(t *testing.T) {
		external, internal, err := DeriveAddresses(nil, 0)
		assert.Error(t, err)
		assert.Equal(t, "", internal)
		assert.Equal(t, "", external)
	})
}

func TestDerivePrivateKeyFromHex(t *testing.T) {
	hdXpriv, _ := compat.GetHDKeyFromExtendedPublicKey(testXpriv)

	t.Run("empty key", func(t *testing.T) {
		_, err := DerivePrivateKeyFromHex(nil, testHash)
		require.ErrorIs(t, err, ErrHDKeyNil)
	})

	t.Run("empty hex", func(t *testing.T) {
		key, err := DerivePrivateKeyFromHex(hdXpriv, "")
		require.NoError(t, err)
		assert.Equal(t, privateKey0, hex.EncodeToString(key.Serialize()))
	})

	t.Run("empty hex", func(t *testing.T) {
		key, err := DerivePrivateKeyFromHex(hdXpriv, testHash)
		require.NoError(t, err)
		assert.Equal(t, privateKeyHash, hex.EncodeToString(key.Serialize()))
	})
}

func TestGetChildNumsFromHex(t *testing.T) {
	t.Run("empty hex", func(t *testing.T) {
		childNums, err := GetChildNumsFromHex("")
		require.NoError(t, err)
		assert.Equal(t, []uint32{}, childNums)
	})

	t.Run("invalid hex", func(t *testing.T) {
		_, err := GetChildNumsFromHex("test")
		require.Error(t, err)
	})

	t.Run("hex ababab", func(t *testing.T) {
		childNums, err := GetChildNumsFromHex("ababab")
		require.NoError(t, err)
		assert.Equal(t, []uint32{11250603}, childNums)
	})

	t.Run("hex testHash", func(t *testing.T) {
		childNums, err := GetChildNumsFromHex(testHash)
		require.NoError(t, err)
		assert.Equal(t, []uint32{
			196136815,  // 8bb0cf6e = 2343620462 - 2147483647
			967933200,  // b9b17d0f = 3115416847 - 2147483647
			2099426390, // 7d22b456
			1897997694, // f121257d = 4045481341 - 2147483647
			1092963872, // c1254e1f = 3240447519 - 2147483647
			23483248,   // 01665370
			1197704170, // 476383ea
			2003694612, // 776df414
		}, childNums)
	})
}

func TestDeriveChildKeyFromHex(t *testing.T) {
	t.Run("xpriv", func(t *testing.T) {
		key, err := compat.GenerateHDKeyFromString(testXpriv)
		require.NoError(t, err)

		var childKey *bip32.ExtendedKey
		childKey, err = DeriveChildKeyFromHex(key, testHash)
		require.NoError(t, err)
		assert.Equal(t, derivedXpriv, childKey.String())
	})

	t.Run("xpub", func(t *testing.T) {
		key, err := compat.GenerateHDKeyFromString(testXpub)
		require.NoError(t, err)

		var childKey *bip32.ExtendedKey
		childKey, err = DeriveChildKeyFromHex(key, testHash)
		require.NoError(t, err)
		assert.Equal(t, derivedXpub, childKey.String())
	})

	t.Run("xpriv => xpub", func(t *testing.T) {
		key, err := compat.GenerateHDKeyFromString(testXpriv)
		require.NoError(t, err)

		var childKey *bip32.ExtendedKey
		childKey, err = DeriveChildKeyFromHex(key, testHash)
		require.NoError(t, err)

		var hdPubKey *bip32.ExtendedKey
		hdPubKey, err = childKey.Neuter()
		require.NoError(t, err)
		assert.Equal(t, derivedXpub, hdPubKey.String())
	})
}

func TestDerivePublicKey(t *testing.T) {
	t.Parallel()

	hdKey, err := ValidateXPub(testXpub)
	require.NoError(t, err)

	t.Run("nil", func(t *testing.T) {
		var pubKey *ec.PublicKey
		pubKey, err = DerivePublicKey(nil, 0, 0)
		assert.ErrorIs(t, err, ErrHDKeyNil)
		assert.Nil(t, pubKey)
	})

	t.Run("derive", func(t *testing.T) {
		var pubKey *ec.PublicKey
		pubKey, err = DerivePublicKey(hdKey, 0, 0)
		require.NoError(t, err)
		assert.Equal(t,
			"03d406421c2733d69a76147c67f8c2194857a2f088299ebf8f1c3790396aa70b4e",
			hex.EncodeToString(pubKey.Compressed()),
		)

		pubKey, err = DerivePublicKey(hdKey, 1, 1)
		require.NoError(t, err)
		assert.Equal(t,
			"0263e4a3696fe4e5136536988169bc3fbec730b912ade7988c57098a47a81a0ae1",
			hex.EncodeToString(pubKey.Compressed()),
		)
	})
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil / empty", func(t *testing.T) {
		assert.False(t, StringInSlice("test", []string{}))
		assert.False(t, StringInSlice("test", nil))
	})

	t.Run("slices", func(t *testing.T) {
		slice := []string{"test", "test1", "test2"}
		assert.True(t, StringInSlice("test", slice))
		assert.True(t, StringInSlice("test1", slice))
		assert.True(t, StringInSlice("test2", slice))
		assert.False(t, StringInSlice("test3", slice))
	})
}

func TestGetTransactionIDFromHex(t *testing.T) {
	t.Parallel()

	t.Run("nil / empty", func(t *testing.T) {
		id, err := GetTransactionIDFromHex("")
		assert.Equal(t, "", id)
		assert.Error(t, err)

		id, err = GetTransactionIDFromHex("test")
		assert.Equal(t, "", id)
		assert.Error(t, err)
	})

	t.Run("tx id from hex", func(t *testing.T) {
		id, err := GetTransactionIDFromHex("020000000165bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10000000006a473044022057b101e9a017cdcc333ef66a4a1e78720ae15adf7d1be9c33abec0fe56bc849d022013daa203095522039fadaba99e567ec3cf8615861d3b7258d5399c9f1f4ace8f412103b9c72aebee5636664b519e5f7264c78614f1e57fa4097ae83a3012a967b1c4b9ffffffff03e0930400000000001976a91413473d21dc9e1fb392f05a028b447b165a052d4d88acf9020000000000001976a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac00000000000000000c006a09446f7457616c6c657400000000")
		assert.Equal(t, "1b52eac9d1eb0adf3ce6a56dee1c4768780b8126e288aca65dd1db32f173b853", id)
		require.NoError(t, err)

		id, err = GetTransactionIDFromHex("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff500307f10a044621b26108fabe6d6d305fc4f4217493d17ce80b94904cab5742b39c5bdfcb21a3b3f216d499e6f9dc000100000000000001a5a2b5eb4c0000db010f2f4d696e696e672d4475746368342f0000000001439cce25000000001976a9141b0cf0b84e60b24cd2bb6391b2d704705657a59188ac00000000")
		assert.Equal(t, "1f123f00823519fd5225c455b06310b534cc7324e2b0cab15108ecaf002a7074", id)
		require.NoError(t, err)
	})
}
