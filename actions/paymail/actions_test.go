package pmail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var externalXPubID = "xpub69PUyEkuD8cqyA9ekUkp3FwaeW1uyLxbwybEy3bmyD7mM6zShsJqfRCv12B43h6KiEiZgF3BFSMnYLsVZr526n37qsqVXkPKYWQ8En2xbi1"
var testPaymail = "paymail@tester.com"

func Test_newPaymailAddress(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		paymail := ""
		_, err := newPaymailAddress(ctx, testXPub, paymail, client.DefaultModelOptions()...)
		require.ErrorIs(t, err, ErrMissingPaymailID)
	})

	t.Run("new paymail address", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		paymailAddress, err := newPaymailAddress(ctx, testXPub, testPaymail, client.DefaultModelOptions()...)
		require.NoError(t, err)

		assert.Equal(t, "paymail", paymailAddress.Alias)
		assert.Equal(t, "tester.com", paymailAddress.Domain)
		assert.Equal(t, testXPubID, paymailAddress.XPubID)
		assert.Equal(t, externalXPubID, paymailAddress.ExternalXPubKey)

		var p2 *PaymailAddress
		p2, err = GetPaymail(ctx, testPaymail, client.DefaultModelOptions()...)
		require.NoError(t, err)

		assert.Equal(t, "paymail", p2.Alias)
		assert.Equal(t, "tester.com", p2.Domain)
		assert.Equal(t, testXPubID, p2.XPubID)
		assert.Equal(t, externalXPubID, p2.ExternalXPubKey)
	})
}

func Test_deletePaymailAddress(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		paymail := ""
		err := deletePaymailAddress(ctx, paymail, client.DefaultModelOptions()...)
		require.ErrorIs(t, err, ErrMissingPaymail)
	})

	t.Run("delete unknown paymail address", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		err := deletePaymailAddress(ctx, testPaymail, client.DefaultModelOptions()...)
		require.ErrorIs(t, err, ErrMissingPaymail)
	})

	t.Run("new paymail address", func(t *testing.T) {
		ctx, client, deferMe := getPaymailClient(t)
		defer deferMe()

		paymailAddress, err := newPaymailAddress(ctx, testXPub, testPaymail, client.DefaultModelOptions()...)
		require.NoError(t, err)

		err = deletePaymailAddress(ctx, testPaymail, client.DefaultModelOptions()...)
		require.NoError(t, err)

		//time.Sleep(1 * time.Second)

		var p2 *PaymailAddress
		p2, err = GetPaymail(ctx, testPaymail, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.Nil(t, p2)

		var p3 *PaymailAddress
		p3, err = GetPaymailByID(ctx, paymailAddress.ID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.Equal(t, testPaymail, p3.Alias)
		require.True(t, p3.DeletedAt.Valid)
	})
}
