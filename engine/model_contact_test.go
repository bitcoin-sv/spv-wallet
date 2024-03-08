package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_contact_validate_success(t *testing.T) {
	t.Run("valid contact", func(t *testing.T) {
		// given
		contact := newContact("Homer Simpson", "homer@springfield.7g", "xpubblablahomer",
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
			expetedError: fmt.Errorf("paymail address failed format validation: definitely not paymail"),
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

		contact := newContact("Homer Simpson", "homer@springfield.7g", "xpubblablahomer",
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

		contact := newContact("Marge Simpson", "Marge@springfield.7g", "xpubblablamarge",
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
