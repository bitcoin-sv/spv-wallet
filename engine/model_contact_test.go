package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func Test_Accept(t *testing.T) {
	t.Run("accept awaiting contact - status changed to <<unconfirmed>>", func(t *testing.T) {
		// given
		sut := newContact("Leto Atreides", "leto@atreides.diune", "pubkey", "xpubid", ContactAwaitAccept)

		// when
		err := sut.Accept()

		// then
		require.NoError(t, err)
		require.Equal(t, ContactNotConfirmed, sut.Status)
	})

	t.Run("accept non-awaiting contact - return error, status has not been changed", func(t *testing.T) {
		// given
		const notAwaitingStatus = ContactNotConfirmed
		sut := newContact("Jessica Atreides", "jess@atreides.diune", "pubkey", "xpubid", notAwaitingStatus)

		// when
		err := sut.Accept()

		// then
		require.Error(t, err)
		require.Equal(t, notAwaitingStatus, sut.Status)
	})
}

func Test_Reject(t *testing.T) {
	t.Run("reject awaiting contact - status changed to <<rejected>>, contact has been marked as deleted", func(t *testing.T) {
		// given
		sut := newContact("Vladimir Harkonnen", "vlad@harkonnen.diune", "pubkey", "xpubid", ContactAwaitAccept)

		// when
		err := sut.Reject()

		// then
		require.NoError(t, err)
		require.Equal(t, ContactRejected, sut.Status)
		require.True(t, sut.DeletedAt.Valid)
	})

	t.Run("reject non-awaiting contact - return error, status has not been changed", func(t *testing.T) {
		// given
		const notAwaitingStatus = ContactNotConfirmed
		sut := newContact("Feyd-Rautha Harkonnen", "frautha@harkonnen.diune", "pubkey", "xpubid", notAwaitingStatus)

		// when
		err := sut.Reject()

		// then
		require.Error(t, err)
		require.Equal(t, notAwaitingStatus, sut.Status)
		require.False(t, sut.DeletedAt.Valid)
	})
}

func Test_Confirm(t *testing.T) {
	t.Run("confirm unconfirmed contact - status changed to <<confirmed>>", func(t *testing.T) {
		// given
		sut := newContact("Thufir Hawat", "hawat@atreides.diune", "pubkey", "xpubid", ContactNotConfirmed)

		// when
		err := sut.Confirm()

		// then
		require.NoError(t, err)
		require.Equal(t, ContactConfirmed, sut.Status)
	})

	t.Run("confirm non-unconfirmed contact - return error, status has not been changed", func(t *testing.T) {
		// given
		const notUncormirmedStatus = ContactAwaitAccept
		sut := newContact("Gurney Halleck", "halleck@atreides.diune", "pubkey", "xpubid", notUncormirmedStatus)

		// when
		err := sut.Confirm()

		// then
		require.Error(t, err)
		require.Equal(t, notUncormirmedStatus, sut.Status)
	})
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

	t.Run("get by paymail for owner xpubid - case insensitive", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		contact := newContact("Homer Simpson", "hOmEr@springfield.com", "xpubblablahomer",
			"fafagasfaufrusfrusfrbsur", ContactNotConfirmed, WithClient(client))

		uppercaseContactPaymail := strings.ToUpper(contact.Paymail)

		err := contact.Save(ctx)
		require.NoError(t, err)

		// when
		result, err := getContact(ctx, uppercaseContactPaymail, contact.OwnerXpubID, WithClient(client))

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

	t.Run("get deleted contact - returns nil", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		contact := newContact("Marge Simpson", "Marge@springfield.com", "xpubblablamarge",
			"fafagasfaufrusfrusfrbsur", ContactNotConfirmed, WithClient(client))

		err := contact.Save(ctx)
		require.NoError(t, err)

		// delete
		contact.DeletedAt.Valid = true
		contact.DeletedAt.Time = time.Now()
		err = contact.Save(ctx)
		require.NoError(t, err)

		// when
		result, err := getContact(ctx, contact.Paymail, contact.OwnerXpubID, WithClient(client))

		// then
		require.NoError(t, err)
		require.Nil(t, result)
	})
}

func Test_getContacts(t *testing.T) {

	t.Run("get by status 'not confirmed'", func(t *testing.T) {
		// given
		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer cleanup()

		xpubID := "xpubid"

		// fullfill db
		saveContactsN(xpubID, ContactAwaitAccept, 10, client)
		saveContactsN(xpubID, ContactNotConfirmed, 13, client)

		conditions := map[string]interface{}{
			contactStatusField: ContactNotConfirmed,
		}

		// when
		contacts, err := getContacts(ctx, xpubID, nil, conditions, nil, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, contacts)
		require.Equal(t, 13, len(contacts))

		for _, c := range contacts {
			require.Equal(t, ContactNotConfirmed, c.Status)
		}

	})

	t.Run("get without conditions - return all", func(t *testing.T) {
		// given
		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer cleanup()

		xpubID := "xpubid"

		// fullfill db
		saveContactsN(xpubID, ContactAwaitAccept, 10, client)
		saveContactsN(xpubID, ContactNotConfirmed, 13, client)

		// when
		contacts, err := getContacts(ctx, xpubID, nil, nil, nil, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, contacts)
		require.Equal(t, 23, len(contacts))

	})

	t.Run("get without conditions - ensure returned only with correct xpubid", func(t *testing.T) {
		// given
		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer cleanup()

		xpubID := "xpubid"

		// fullfill db
		saveContactsN(xpubID, ContactAwaitAccept, 10, client)
		saveContactsN("other-xpub", ContactNotConfirmed, 13, client)

		// when
		contacts, err := getContacts(ctx, xpubID, nil, nil, nil, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, contacts)
		require.Equal(t, 10, len(contacts))

	})

	t.Run("get without conditions - ensure returned without deleted", func(t *testing.T) {
		// given
		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer cleanup()

		xpubID := "xpubid"

		// fullfill db
		saveContactsN(xpubID, ContactAwaitAccept, 10, client)
		saveContactsDeletedN(xpubID, ContactNotConfirmed, 13, client)

		// when
		contacts, err := getContacts(ctx, xpubID, nil, nil, nil, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, contacts)
		require.Equal(t, 10, len(contacts))

	})
}

func saveContactsN(xpubID string, status ContactStatus, n int, c ClientInterface) {
	for i := 0; i < n; i++ {
		e := newContact(fmt.Sprintf("%s%d", status, i), fmt.Sprintf("%s%d@t.com", status, i), "pubkey", xpubID, status, c.DefaultModelOptions()...)
		if err := e.Save(context.Background()); err != nil {
			panic(err)
		}
	}
}

func saveContactsDeletedN(xpubID string, status ContactStatus, n int, c ClientInterface) {
	for i := 0; i < n; i++ {
		e := newContact(fmt.Sprintf("%s%d", status, i), fmt.Sprintf("%s%d@t.com", status, i), "pubkey", xpubID, status, c.DefaultModelOptions()...)
		e.DeletedAt.Valid = true
		e.DeletedAt.Time = time.Now()
		if err := e.Save(context.Background()); err != nil {
			panic(err)
		}
	}
}
