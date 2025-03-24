package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contactFixture struct {
	engine        engine.ClientInterface
	user          fixtures.User
	t             testing.TB
	assert        *assert.Assertions
	require       *require.Assertions
	paymailClient *paymailmock.PaymailClientMock
}

const (
	examplePubKey1 = "0303b140f2a317069858db56e9d09aed87d8ffab2a8f70960eceec10f0a67d5053"
	examplePubKey2 = "03736d66a37481d1c3784bcbf60997c0fb2c4e2691e63129ff33074a4580b445ef"
)

func (f *contactFixture) HasContactTo(userB fixtures.User) *contactsmodels.Contact {
	newContact := contactsmodels.NewContact{
		UserID:            f.user.ID(),
		FullName:          userB.DefaultPaymail().PublicName(),
		NewContactPaymail: userB.DefaultPaymail().String(),
		NewContactPubKey:  examplePubKey1,
		Status:            contactsmodels.ContactNotConfirmed,
	}

	contact, err := f.engine.Repositories().Contacts.Create(context.Background(), newContact)
	f.require.NoError(err)

	return contact
}

func (f *contactFixture) HasConfirmedContactTo(userB fixtures.User) *contactsmodels.Contact {
	newContact := contactsmodels.NewContact{
		UserID:            f.user.ID(),
		FullName:          userB.DefaultPaymail().PublicName(),
		NewContactPaymail: userB.DefaultPaymail().String(),
		NewContactPubKey:  examplePubKey1,
		Status:            contactsmodels.ContactConfirmed,
	}

	contact, err := f.engine.Repositories().Contacts.Create(context.Background(), newContact)
	f.require.NoError(err)

	return contact
}

func (f *contactFixture) HasRejectedContactTo(userB fixtures.User) *contactsmodels.Contact {
	newContact := contactsmodels.NewContact{
		UserID:            f.user.ID(),
		FullName:          userB.DefaultPaymail().PublicName(),
		NewContactPaymail: userB.DefaultPaymail().String(),
		NewContactPubKey:  examplePubKey1,
		Status:            contactsmodels.ContactRejected,
	}

	contact, err := f.engine.Repositories().Contacts.Create(context.Background(), newContact)
	f.require.NoError(err)

	return contact
}

func (f *contactFixture) HasAwaitingContactTo(userB fixtures.User) *contactsmodels.Contact {
	newContact := contactsmodels.NewContact{
		UserID:            f.user.ID(),
		FullName:          userB.DefaultPaymail().PublicName(),
		NewContactPaymail: userB.DefaultPaymail().String(),
		NewContactPubKey:  examplePubKey1,
		Status:            contactsmodels.ContactAwaitAccept,
	}

	contact, err := f.engine.Repositories().Contacts.Create(context.Background(), newContact)
	f.require.NoError(err)

	return contact
}
