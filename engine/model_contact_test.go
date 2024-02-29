package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fullName     = "John Doe"
	paymailTest  = "test@paymail.com"
	senderPubKey = "senderPubKey"
)

func Test_newContact(t *testing.T) {
	t.Run("valid full_name and paymail, pubKey is not empty", func(t *testing.T) {
		contact, err := newContact(fullName, paymailTest, senderPubKey)
		require.NoError(t, err)
		require.NotNil(t, contact)

		assert.Equal(t, fullName, contact.FullName)
		assert.Equal(t, paymailTest, contact.Paymail)
		assert.Equal(t, utils.Hash(senderPubKey), contact.XpubID)
		assert.Equal(t, utils.Hash(contact.XpubID+paymailTest), contact.ID)
	})

	t.Run("empty full_name", func(t *testing.T) {
		contact, err := newContact("", paymailTest, senderPubKey)
		require.Errorf(t, err, ErrEmptyContactFullName.Error())
		require.EqualError(t, err, ErrEmptyContactFullName.Error())

		require.Nil(t, contact)
	})

	t.Run("empty paymail", func(t *testing.T) {
		contact, err := newContact(fullName, "", senderPubKey)

		require.Nil(t, contact)
		require.ErrorContains(t, err, "paymail address failed format validation")
	})

	t.Run("invalid paymail", func(t *testing.T) {
		contact, err := newContact(fullName, "testata", senderPubKey)

		require.Nil(t, contact)
		require.ErrorContains(t, err, "paymail address failed format validation")
	})

	t.Run("empty pubKey", func(t *testing.T) {

		contact, err := newContact(fullName, paymailTest, "")

		require.Nil(t, contact)
		require.ErrorContains(t, err, ErrEmptyContactPubKey.Error())
	})
}

func Test_getContact(t *testing.T) {
	t.Run("fullName, paymail and pubKey are valid", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, paymailTest, senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.NoError(t, err)
	})

	t.Run("empty full_name", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, "", paymailTest, senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.Error(t, err)
	})

	t.Run("empty paymail", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, "", senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.Error(t, err)
	})

	t.Run("invalid paymail", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, "tests", senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.Error(t, err)
	})

	t.Run("empty pubKey", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, paymailTest, "", client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.Error(t, err)
	})
}
