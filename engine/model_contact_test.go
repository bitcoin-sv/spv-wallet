package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fullName     = "John Doe"
	paymailTest  = "test@paymail.com"
	senderPubKey = "senderPubKey"
	xPubID       = "62910a1ecbc7728afad563ab3f8aa70568ed934d1e0383cb1bbbfb1bc8f2afe5"
)

func Test_newContact(t *testing.T) {
	t.Run("valid full_name and paymail, pubKey is not empty", func(t *testing.T) {
		contact, err := newContact(fullName, paymailTest, senderPubKey)
		require.NoError(t, err)
		require.NotNil(t, contact)
		assert.Equal(t, fullName, contact.FullName)
		assert.Equal(t, paymailTest, contact.Paymail)
		assert.Equal(t, utils.Hash(senderPubKey), contact.XpubID)
		assert.Equal(t, utils.Hash(senderPubKey+paymailTest), contact.ID)
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
		require.EqualError(t, err, ErrEmptyContactPaymail.Error())
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
		require.NoError(t, err)
	})

	t.Run("empty paymail", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, "", senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.NoError(t, err)
	})

	t.Run("invalid paymail", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, "tests", senderPubKey, client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.NoError(t, err)
	})

	t.Run("empty pubKey", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		contact, err := getContact(ctx, fullName, paymailTest, "", client.DefaultModelOptions()...)

		require.Nil(t, contact)
		require.NoError(t, err)
	})
}

func Test_getContacts(t *testing.T) {
	t.Run("status 'not confirmed'", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		var metadata *Metadata

		dbConditions := make(map[string]interface{})

		var queryParams *datastore.QueryParams

		(dbConditions)[xPubIDField] = xPubID
		(dbConditions)[contactStatus] = notConfirmed

		contacts, err := getContacts(ctx, metadata, &dbConditions, queryParams, client.DefaultModelOptions()...)

		require.NoError(t, err)
		assert.NotNil(t, contacts)
	})

	t.Run("status 'confirmed'", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		var metadata *Metadata

		dbConditions := make(map[string]interface{})

		var queryParams *datastore.QueryParams

		(dbConditions)[xPubIDField] = xPubID
		(dbConditions)[contactStatus] = confirmed

		contacts, err := getContacts(ctx, metadata, &dbConditions, queryParams, client.DefaultModelOptions()...)

		require.NoError(t, err)
		assert.Equal(t, 0, len(contacts))
	})

	t.Run("status 'awaiting acceptance'", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		var metadata *Metadata

		dbConditions := make(map[string]interface{})

		var queryParams *datastore.QueryParams

		(dbConditions)[xPubIDField] = xPubID
		(dbConditions)[contactStatus] = awaitingAcceptance

		contacts, err := getContacts(ctx, metadata, &dbConditions, queryParams, client.DefaultModelOptions()...)

		require.NoError(t, err)
		assert.Equal(t, 0, len(contacts))
	})
}
