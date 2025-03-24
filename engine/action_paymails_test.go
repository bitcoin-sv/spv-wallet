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
	testPaymail    = "paymail@tester.com"
	testPublicName = "Public Name"
)

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
			require.Nil(t, p3)
		})
	}
}
