package engine

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	externalXPubID = "xpub6BYQwVS1dNAYoVVcN7cASn5qv4E7QEX4dCazSkJwAk89w32pAfbJqDLfibRW4ywEBwus54uCD8PwLYzahiyuMbyLujCT2oQD5z6QobaNyN1"
	testAvatar     = "https://i.imgur.com/MYSVX44.png"
	testAvatar2    = "https://i.imgur.com/cBJKPDh.png"
	testPaymail    = "paymail@tester.com"
	testPublicName = "Public Name"
)

// TestClient_NewPaymailAddress will test the method NewPaymailAddress()
func (ts *EmbeddedDBTestSuite) TestClient_NewPaymailAddress() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - empty address", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			// Create xPub (required to add a paymail address)
			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.NotNil(t, xPub)
			require.NoError(t, err)

			var paymailAddress *PaymailAddress
			paymailAddress, err = tc.client.NewPaymailAddress(tc.ctx, testXPub, "", testPublicName, testAvatar, tc.client.DefaultModelOptions()...)
			require.ErrorIs(t, err, spverrors.ErrMissingPaymailAddress)
			require.Nil(t, paymailAddress)
		})

		ts.T().Run(testCase.name+" - new paymail address", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			opts := tc.client.DefaultModelOptions()

			// Create xPub (required to add a paymail address)
			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NotNil(t, xPub)
			require.NoError(t, err)

			var paymailAddress *PaymailAddress
			paymailAddress, err = tc.client.NewPaymailAddress(tc.ctx, xPub.RawXpub(), testPaymail, testPublicName, testAvatar, opts...)
			require.NoError(t, err)
			require.NotNil(t, paymailAddress)

			assert.Equal(t, "paymail", paymailAddress.Alias)
			assert.Equal(t, "tester.com", paymailAddress.Domain)
			assert.Equal(t, testAvatar, paymailAddress.Avatar)
			assert.Equal(t, testPublicName, paymailAddress.PublicName)
			assert.Equal(t, testXPubID, paymailAddress.XpubID)
			assert.Equal(t, externalXPubID, paymailAddress.ExternalXpubKey)

			var p2 *PaymailAddress
			p2, err = getPaymailAddress(tc.ctx, testPaymail, opts...)
			require.NoError(t, err)
			require.NotNil(t, p2)

			assert.Equal(t, "paymail", p2.Alias)
			assert.Equal(t, "tester.com", p2.Domain)
			assert.Equal(t, testAvatar, p2.Avatar)
			assert.Equal(t, testPublicName, p2.PublicName)
			assert.Equal(t, testXPubID, p2.XpubID)
			assert.Equal(t, externalXPubID, p2.ExternalXpubKey)
		})
	}
}

// Test_DeletePaymailAddress will test the method DeletePaymailAddress()
func (ts *EmbeddedDBTestSuite) Test_DeletePaymailAddress() {
	for _, testCase := range dbTestCases {

		ts.T().Run(testCase.name+" - empty", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			paymail := ""
			err := tc.client.DeletePaymailAddress(tc.ctx, paymail, tc.client.DefaultModelOptions()...)
			require.ErrorIs(t, err, spverrors.ErrCouldNotFindPaymail)
		})

		ts.T().Run(testCase.name+" - delete unknown paymail address", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			err := tc.client.DeletePaymailAddress(tc.ctx, testPaymail, tc.client.DefaultModelOptions()...)
			require.ErrorIs(t, err, spverrors.ErrCouldNotFindPaymail)
		})

		ts.T().Run(testCase.name+" - delete paymail address", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)
			opts := tc.client.DefaultModelOptions()

			// Create xPub (required to add a paymail address)
			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NotNil(t, xPub)
			require.NoError(t, err)

			var paymailAddress *PaymailAddress
			paymailAddress, err = tc.client.NewPaymailAddress(tc.ctx, testXPub, testPaymail, testPublicName, testAvatar, opts...)
			require.NoError(t, err)
			require.NotNil(t, paymailAddress)

			err = tc.client.DeletePaymailAddress(tc.ctx, testPaymail, opts...)
			require.NoError(t, err)

			var p2 *PaymailAddress
			p2, err = getPaymailAddress(tc.ctx, testPaymail, opts...)
			require.NoError(t, err)
			require.Nil(t, p2)

			var p3 *PaymailAddress
			p3, err = getPaymailAddressByID(tc.ctx, paymailAddress.ID, opts...)
			require.NoError(t, err)
			require.NotNil(t, p3)
			require.Equal(t, testPaymail, p3.Alias)
			require.True(t, p3.DeletedAt.Valid)
		})
	}
}

// TestClient_UpdatePaymailAddressMetadata will test the method UpdatePaymailAddressMetadata()
func (ts *EmbeddedDBTestSuite) TestClient_UpdatePaymailAddressMetadata() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			metadata := Metadata{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
				"test-key-3": "test-value-3",
			}
			opts := tc.client.DefaultModelOptions()
			opts = append(opts, WithMetadatas(metadata))

			// Create xPub (required to add a paymail address)
			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NotNil(t, xPub)
			require.NoError(t, err)

			var paymailAddress *PaymailAddress
			paymailAddress, err = tc.client.NewPaymailAddress(tc.ctx, testXPub, testPaymail, testPublicName, testAvatar, opts...)
			require.NoError(t, err)
			require.NotNil(t, paymailAddress)

			paymailAddress, err = tc.client.UpdatePaymailAddressMetadata(tc.ctx, testPaymail, Metadata{"test-key-new": "new-value"}, opts...)
			require.NoError(t, err)
			assert.Len(t, paymailAddress.Metadata, 4)
			assert.Equal(t, "new-value", paymailAddress.Metadata["test-key-new"])

			paymailAddress, err = tc.client.UpdatePaymailAddressMetadata(tc.ctx, testPaymail, Metadata{
				"test-key-new-2": "new-value-2",
				"test-key-1":     nil,
				"test-key-2":     nil,
				"test-key-3":     nil,
			}, opts...)
			require.NoError(t, err)
			assert.Len(t, paymailAddress.Metadata, 2)
			assert.Equal(t, "new-value", paymailAddress.Metadata["test-key-new"])
			assert.Equal(t, "new-value-2", paymailAddress.Metadata["test-key-new-2"])

			var p2 *PaymailAddress
			p2, err = getPaymailAddress(tc.ctx, testPaymail, opts...)
			require.NoError(t, err)
			require.NotNil(t, p2)
			assert.Len(t, paymailAddress.Metadata, 2)
		})
	}
}

// TestClient_UpdatePaymailAddress will test the method UpdatePaymailAddress()
func (ts *EmbeddedDBTestSuite) TestClient_UpdatePaymailAddress() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			opts := tc.client.DefaultModelOptions()

			// Create xPub (required to add a paymail address)
			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NotNil(t, xPub)
			require.NoError(t, err)

			var paymailAddress *PaymailAddress
			paymailAddress, err = tc.client.NewPaymailAddress(tc.ctx, testXPub, testPaymail, testPublicName, testAvatar, opts...)
			require.NoError(t, err)
			require.NotNil(t, paymailAddress)
			assert.Equal(t, testPublicName, paymailAddress.PublicName)
			assert.Equal(t, testAvatar, paymailAddress.Avatar)

			paymailAddress, err = tc.client.UpdatePaymailAddress(tc.ctx, testPaymail, testPublicName+"2", testAvatar2, opts...)
			require.NoError(t, err)

			assert.Equal(t, testPublicName+"2", paymailAddress.PublicName)
			assert.Equal(t, testAvatar2, paymailAddress.Avatar)

			var p2 *PaymailAddress
			p2, err = getPaymailAddress(tc.ctx, testPaymail, tc.client.DefaultModelOptions()...)
			require.NoError(t, err)
			require.NotNil(t, p2)
			assert.Equal(t, testPublicName+"2", p2.PublicName)
			assert.Equal(t, testAvatar2, p2.Avatar)
		})
	}
}
