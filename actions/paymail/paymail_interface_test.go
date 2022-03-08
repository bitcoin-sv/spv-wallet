package pmail

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/tester"
	"github.com/BuxOrg/bux/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-paymail"
	"github.com/tonicpow/go-paymail/server"
)

const (
	alias       = "paymail"
	domain      = "tester.com"
	fullPaymail = "paymail@tester.com"
)

var (
	// testTxHex = "020000000165bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10000000006a473044022057b101e9a017cdcc333ef66a4a1e78720ae15adf7d1be9c33abec0fe56bc849d022013daa203095522039fadaba99e567ec3cf8615861d3b7258d5399c9f1f4ace8f412103b9c72aebee5636664b519e5f7264c78614f1e57fa4097ae83a3012a967b1c4b9ffffffff03e0930400000000001976a91413473d21dc9e1fb392f05a028b447b165a052d4d88acf9020000000000001976a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac00000000000000000c006a09446f7457616c6c657400000000"
	// testTxID  = "1b52eac9d1eb0adf3ce6a56dee1c4768780b8126e288aca65dd1db32f173b853"
	testXPub   = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	testXPubID = "1a0b10d4eda0636aae1709e7e7080485a4d99af3ca2962c6e677cf5b53d8ab8c"
)

func TestPaymailInterface(t *testing.T) {
	t.Parallel()

	t.Run("GetPaymailByAlias", func(t *testing.T) {
		ctx, client, deferMe, xPub, paymailModelService, externalXPubKey, external := initPaymailTesting(t)
		defer deferMe()

		paymailAddress, err := paymailModelService.GetPaymailByAlias(ctx, alias, domain, nil)
		require.NoError(t, err)
		assert.IsType(t, paymail.AddressInformation{}, *paymailAddress)
		assert.Equal(t, alias, paymailAddress.Alias)
		assert.Equal(t, domain, paymailAddress.Domain)
		assert.Equal(t, externalXPubKey, paymailAddress.PubKey)
		assert.Equal(t, external, paymailAddress.LastAddress)
		assert.Equal(t, "Tester", paymailAddress.Name)

		destination := checkCreatedDestination(ctx, t, client, xPub, external, "GetPaymailByAlias")
		assert.Equal(t, "GetPaymailByAlias", destination.Metadata[paymailRequestField])
	})

	t.Run("GetPaymailByAlias with metadata", func(t *testing.T) {
		ctx, client, deferMe, xPub, paymailModelService, externalXPubKey, external := initPaymailTesting(t)
		defer deferMe()

		metadata := &server.RequestMetadata{
			UserAgent: "test-user-agent",
			Note:      "test-note",
		}
		paymailAddress, err := paymailModelService.GetPaymailByAlias(ctx, alias, domain, metadata)
		require.NoError(t, err)
		assert.IsType(t, paymail.AddressInformation{}, *paymailAddress)
		assert.Equal(t, alias, paymailAddress.Alias)
		assert.Equal(t, domain, paymailAddress.Domain)
		assert.Equal(t, externalXPubKey, paymailAddress.PubKey)
		assert.Equal(t, external, paymailAddress.LastAddress)
		assert.Equal(t, "Tester", paymailAddress.Name)

		destination := checkCreatedDestination(ctx, t, client, xPub, external, "GetPaymailByAlias")
		assert.Equal(t, "GetPaymailByAlias", destination.Metadata[paymailRequestField])
		assert.Equal(t, "test-user-agent", destination.Metadata["user-agent"])
		assert.Equal(t, "test-note", destination.Metadata["note"])
	})

	t.Run("CreateAddressResolutionResponse", func(t *testing.T) {
		ctx, client, deferMe, xPub, paymailModelService, _, external := initPaymailTesting(t)
		defer deferMe()

		resolutionInformation, err := paymailModelService.CreateAddressResolutionResponse(ctx, alias, domain, false, nil)
		require.NoError(t, err)
		assert.IsType(t, paymail.ResolutionPayload{}, *resolutionInformation)
		assert.Equal(t, external, resolutionInformation.Address)

		destination := checkCreatedDestination(ctx, t, client, xPub, external, "CreateAddressResolutionResponse")
		assert.Equal(t, destination.LockingScript, resolutionInformation.Output)
		assert.Equal(t, "CreateAddressResolutionResponse", destination.Metadata[paymailRequestField])
	})

	t.Run("CreateP2PDestinationResponse", func(t *testing.T) {
		ctx, client, deferMe, xPub, paymailModelService, _, external := initPaymailTesting(t)
		defer deferMe()

		paymentDestinationInformation, err := paymailModelService.CreateP2PDestinationResponse(ctx, alias, domain, 12000, nil)
		require.NoError(t, err)
		assert.IsType(t, paymail.PaymentDestinationPayload{}, *paymentDestinationInformation)

		destination := checkCreatedDestination(ctx, t, client, xPub, external, "CreateP2PDestinationResponse")

		assert.Equal(t, 1, len(paymentDestinationInformation.Outputs))
		assert.Equal(t, destination.Address, paymentDestinationInformation.Outputs[0].Address)
		assert.Equal(t, uint64(12000), paymentDestinationInformation.Outputs[0].Satoshis)
		assert.Equal(t, destination.LockingScript, paymentDestinationInformation.Outputs[0].Script)
		assert.Equal(t, destination.Metadata[bux.ReferenceIDField], paymentDestinationInformation.Reference)
	})

	// todo: fix this test! (add missing tests)
	/*
		t.Run("RecordTransaction", func(t *testing.T) {
			ctx, _, _, paymailModelService, _, _ := initPaymailTesting(t)

			p2pTx := &paymail.P2PTransaction{
				Hex: testTxHex,
				MetaData: &paymail.P2PMetaData{
					Note:      "test note",
					PubKey:    "some pub key",
					Sender:    "I am the sender",
					Signature: "some signature",
				},
				Reference: "myReferenceID",
			}

			p2PTransactionResponse, err := paymailModelService.RecordTransaction(ctx, p2pTx, nil)
			require.NoError(t, err)
			assert.IsType(t, paymail.P2PTransactionResponse{}, *p2PTransactionResponse)
			assert.IsType(t, testTxID, p2PTransactionResponse.TxID)
		})
	*/
}

func checkCreatedDestination(ctx context.Context, t *testing.T, client bux.ClientInterface, xPub *bux.Xpub,
	external, paymailMetaSignature string) *bux.Destination {

	// check that the destination was created properly
	destination, err := client.GetDestinationByAddress(ctx, testXPubID, external)
	require.NoError(t, err)

	assert.IsType(t, bux.Destination{}, *destination)
	assert.Equal(t, xPub.ID, destination.XpubID)
	assert.Equal(t, uint32(0), destination.Chain)
	assert.Equal(t, uint32(0), destination.Num)
	assert.Equal(t, external, destination.Address)
	assert.Equal(t, paymailMetaSignature, destination.Metadata[paymailRequestField])
	assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)

	return destination
}

func initPaymailTesting(t *testing.T) (context.Context, bux.ClientInterface, func(), *bux.Xpub,
	*PaymailInterface, string, string) {
	ctx, client, deferMe := getPaymailClient(t)

	xPub, err := client.NewXpub(ctx, testXPub)
	require.NoError(t, err)
	require.NotNil(t, xPub)
	assert.IsType(t, bux.Xpub{}, *xPub)

	var hdKey *bip32.ExtendedKey
	hdKey, err = utils.ValidateXPub(testXPub)
	require.NoError(t, err)
	require.NotNil(t, hdKey)

	// derive the first child for the fullPaymail xPub
	var paymailKey *bip32.ExtendedKey
	paymailKey, err = bitcoin.GetHDKeyChild(hdKey, utils.ChainExternal)
	require.NoError(t, err)
	require.NotNil(t, paymailKey)

	// derive the second child for the address / pubKey
	var externalPaymailKey *bip32.ExtendedKey
	externalPaymailKey, err = bitcoin.GetHDKeyChild(paymailKey, 0)
	require.NoError(t, err)

	var externalPaymailXPub *bec.PublicKey
	externalPaymailXPub, err = externalPaymailKey.ECPubKey()
	require.NoError(t, err)

	externalXPubKey := hex.EncodeToString(externalPaymailXPub.SerialiseCompressed())

	// todo: this needs a function or cleanup?
	savePaymailAddress := &bux.PaymailAddress{
		Alias:           alias,
		Avatar:          "img url",
		Domain:          domain,
		ExternalXPubKey: paymailKey.String(),
		ID:              utils.Hash(fullPaymail),
		Model:           *bux.NewBaseModel(bux.ModelPaymail, client.DefaultModelOptions()...),
		Username:        "Tester",
		XPubID:          xPub.ID,
	}
	err = savePaymailAddress.Save(ctx)
	require.NoError(t, err)

	paymailModelService := new(PaymailInterface)
	paymailModelService.client = client

	c := []byte("{\"paymail_server\": {\n    \"enabled\": true,\n    \"domains\": [\n      \"localhost\"\n    ],\n    \"sender_validation_enabled\": false\n  }}")
	err = json.Unmarshal(c, &paymailModelService.appConfig)
	require.NoError(t, err)

	// expected address, derived from the full xPub
	var external string
	external, _, err = utils.DeriveAddresses(hdKey, 0)
	require.NoError(t, err)

	return ctx, client, deferMe, xPub, paymailModelService, externalXPubKey, external
}

func getPaymailClient(t *testing.T) (context.Context, bux.ClientInterface, func()) {
	ctx := context.Background()
	client, err := bux.NewClient(ctx,
		bux.WithSQLite(tester.SQLiteTestConfig(t, true, false)),
		bux.WithRistretto(cachestore.DefaultRistrettoConfig()),
		bux.WithTaskQ(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix(t)+"_queue"), taskmanager.FactoryMemory),
		bux.WithDebugging(),
		bux.WithAutoMigrate(append(bux.BaseModels, &bux.PaymailAddress{})...),
	)
	require.NoError(t, err)

	// Create a defer function
	f := func() {
		_ = client.Close(ctx)
	}

	return ctx, client, f
}
