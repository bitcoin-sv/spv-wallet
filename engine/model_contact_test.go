package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const xPubID = "62910a1ecbc7728afad563ab3f8aa70568ed934d1e0383cb1bbbfb1bc8f2afe5"

func Test_contact_validate_success(t *testing.T) {
	t.Run("valid contact", func(t *testing.T) {
		// given
		contact := newContact("Homer Simpson", "homer@springfield.gg", "xpubblablahomer",
			"fafagasfaufrusfrusfrbsur", ContactNotConfirmed)

		// when
		err := contact.validate()

		// then
		require.NoError(t, err)
		require.NotNil(t, contact)
	})
}

func Test_contact_validate_returns_error(t *testing.T) {
	tcs := []struct {
		name         string
		contact      *Contact
		expetedError error
	}{
		{
			name:         "empty full name",
			contact:      newContact("", "donot@know.who", "xpubblablablabla", "ownerspbubid", ContactNotConfirmed),
			expetedError: ErrMissingContactFullName,
		},
		{
			name:         "empty paymail",
			contact:      newContact("Homer Simpson", "", "xpubblablahomer", "ownerspbubid", ContactNotConfirmed),
			expetedError: errors.New("paymail address failed format validation: "),
		},
		{
			name:         "invalid paymail",
			contact:      newContact("Marge Simpson", "definitely not paymail", "xpubblablamarge", "ownerspbubid", ContactNotConfirmed),
			expetedError: fmt.Errorf("paymail address failed format validation: definitelynotpaymail"),
		},
		{
			name:         "empty pubKey",
			contact:      newContact("Bart Simpson", "bart@springfield.com", "", "ownerspbubid", ContactNotConfirmed),
			expetedError: ErrMissingContactXPubKey,
		},
		{
			name:         "no owner id",
			contact:      newContact("Lisa Simpson", "lisa@springfield.com", "xpubblablalisa", "", ContactNotConfirmed),
			expetedError: ErrMissingContactOwnerXPubId,
		},
		{
			name:         "no status",
			contact:      newContact("Margaret Simpson", "maggie@springfield.com", "xpubblablamaggie", "ownerspbubid", ""),
			expetedError: ErrMissingContactStatus,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// given
			contact := tc.contact

			// when
			err := contact.validate()

			// then
			require.EqualError(t, err, tc.expetedError.Error())
		})
	}
}

func Test_getContact(t *testing.T) {
	t.Run("get by paymail for owner xpubid", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		contact := newContact("Homer Simpson", "homer@springfield.com", "xpubblablahomer",
			"fafagasfaufrusfrusfrbsur", ContactNotConfirmed, WithClient(client))

		err := contact.Save(ctx)
		require.NoError(t, err)

		// when
		result, err := getContact(ctx, contact.Paymail, contact.OwnerXpubID, WithClient(client))

		// then
		require.NoError(t, err)
		require.NotNil(t, result)

		require.Equal(t, contact.ID, result.ID)
		require.Equal(t, contact.OwnerXpubID, result.OwnerXpubID)
		require.Equal(t, contact.FullName, result.FullName)
		require.Equal(t, contact.Paymail, result.Paymail)
		require.Equal(t, contact.PubKey, result.PubKey)
		require.Equal(t, contact.Status, result.Status)
	})

	t.Run("get by paymail for not matching owner xpubid - returns nil", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		contact := newContact("Marge Simpson", "Marge@springfield.com", "xpubblablamarge",
			"fafagasfaufrusfrusfrbsur", ContactNotConfirmed, WithClient(client))

		err := contact.Save(ctx)
		require.NoError(t, err)

		// when
		result, err := getContact(ctx, contact.Paymail, "not owner xpubid", WithClient(client))

		// then
		require.NoError(t, err)
		require.Nil(t, result)
	})
}

func Test_getContacts(t *testing.T) {
	t.Run("status 'not confirmed'", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		var metadata *Metadata

		dbConditions := map[string]interface{}{
			xPubIDField:   xPubID,
			contactStatus: ContactNotConfirmed,
		}

		var queryParams *datastore.QueryParams

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
		(dbConditions)[contactStatus] = ContactConfirmed

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
		(dbConditions)[contactStatus] = ContactAwaitAccept

		contacts, err := getContacts(ctx, metadata, &dbConditions, queryParams, client.DefaultModelOptions()...)

		require.NoError(t, err)
		assert.Equal(t, 0, len(contacts))

	})
}
