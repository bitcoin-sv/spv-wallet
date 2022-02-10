package pmail

import (
	"context"
	"testing"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/tester"
	"github.com/BuxOrg/bux/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bip32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// testTxHex = "020000000165bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10000000006a473044022057b101e9a017cdcc333ef66a4a1e78720ae15adf7d1be9c33abec0fe56bc849d022013daa203095522039fadaba99e567ec3cf8615861d3b7258d5399c9f1f4ace8f412103b9c72aebee5636664b519e5f7264c78614f1e57fa4097ae83a3012a967b1c4b9ffffffff03e0930400000000001976a91413473d21dc9e1fb392f05a028b447b165a052d4d88acf9020000000000001976a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac00000000000000000c006a09446f7457616c6c657400000000"
	// testTxID  = "1b52eac9d1eb0adf3ce6a56dee1c4768780b8126e288aca65dd1db32f173b853"
	testXPub = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
)

// todo: refactor, cleanup, test name, add more tests etc

func TestNewPaymail(t *testing.T) {

	t.Run("paymail basic test", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		paymail := "paymail@tester.com"
		xPubID := utils.Hash(testXPub)

		hdKey, err := utils.ValidateXPub(testXPub)
		require.NoError(t, err)
		var paymailKey *bip32.ExtendedKey
		paymailKey, err = bitcoin.GetHDKeyChild(hdKey, utils.ChainExternal)
		require.NoError(t, err)
		paymailXPub := paymailKey.String()

		p := NewPaymail(
			paymail,
			bux.WithClient(client),
		)
		p.Username = "Tester"
		p.Avatar = "img url"
		p.XPubID = xPubID
		p.ExternalXPubKey = paymailXPub
		err = p.Save(ctx)
		require.NoError(t, err)

		p2 := NewPaymail(paymail, client.DefaultModelOptions()...)
		conditions := map[string]interface{}{
			"alias":  p.Alias,
			"domain": p.Domain,
		}
		err = bux.Get(ctx, p2, conditions, false, 0)
		require.NoError(t, err)
		assert.Equal(t, paymail, p2.Alias+"@"+p2.Domain)
		assert.Equal(t, "Tester", p2.Username)
		assert.Equal(t, "img url", p2.Avatar)
		assert.Equal(t, xPubID, p2.XPubID)
		assert.Equal(t, paymailXPub, p2.ExternalXPubKey)
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

		childKeyChain0, _ := bitcoin.GetHDKeyChild(hdKey, 0)
		childKeyChain01, _ := bitcoin.GetHDKeyChild(childKeyChain0, 1)
		key0, _ := childKeyChain01.ECPubKey()
		address0, _ := bitcoin.GetAddressFromPubKey(key0, true)
		assert.Equal(t, addressExternal, address0.AddressString)

		childKeyChain1, _ := bitcoin.GetHDKeyChild(hdKey, 1)
		childKeyChain11, _ := bitcoin.GetHDKeyChild(childKeyChain1, 1)
		key1, _ := childKeyChain11.ECPubKey()
		address1, _ := bitcoin.GetAddressFromPubKey(key1, true)
		assert.Equal(t, addressInternal, address1.AddressString)
	})
}

func getPaymailClient(t *testing.T) (context.Context, bux.ClientInterface, func()) {
	ctx := context.Background()
	client, err := bux.NewClient(ctx,
		bux.WithSQLite(tester.SQLiteTestConfig(t, true, false)),
		bux.WithRistretto(cachestore.DefaultRistrettoConfig()),
		bux.WithTaskQ(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix(t)+"_queue"), taskmanager.FactoryMemory),
		bux.WithDebugging(),
		bux.WithAutoMigrate(append(bux.BaseModels, &PaymailAddress{})...),
	)
	require.NoError(t, err)

	// Create a defer function
	f := func() {
		_ = client.Close(ctx)
	}

	return ctx, client, f
}
