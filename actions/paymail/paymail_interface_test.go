package pmail

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/BuxOrg/bux"
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
	savePaymailAddress := &PaymailAddress{
		Alias:           alias,
		Avatar:          "img url",
		Domain:          domain,
		ExternalXPubKey: paymailKey.String(),
		ID:              utils.Hash(fullPaymail),
		Model:           *bux.NewBaseModel(ModelPaymail, client.DefaultModelOptions()...),
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
