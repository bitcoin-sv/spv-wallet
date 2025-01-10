package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type assertion struct {
	t       testing.TB
	require *require.Assertions
	assert  *assert.Assertions
	err     error
	contact *engine.Contact
}

type ContactFailureAssertion interface {
	ThatIs(expectedErr error) ContactFailureAssertion
}

type ContactSuccessAssertion interface {
	WithStatus(status engine.ContactStatus) ContactSuccessAssertion
	AsNotConfirmed() ContactSuccessAssertion
	WithFullName(fullName string) ContactSuccessAssertion
	ForUser(user fixtures.User) ContactSuccessAssertion
	ToCounterparty(user fixtures.User) ContactSuccessAssertion
	WithPubKey(pki string) ContactSuccessAssertion
}

type ContactErrorAssertion interface {
	WithError(err error) ContactFailureAssertion
	WithNoError(err error) ContactSuccessAssertion
}
type ContactAssertion interface {
	Created(contact *engine.Contact) ContactErrorAssertion
}

func then(t testing.TB) ContactAssertion {
	return &assertion{
		t:       t,
		require: require.New(t),
		assert:  assert.New(t),
	}
}

func (a *assertion) WithError(err error) ContactFailureAssertion {
	a.require.Nil(a.contact, "unexpected response")
	a.require.Error(err, "error expected")
	a.err = err
	return a
}

func (a *assertion) WithNoError(err error) ContactSuccessAssertion {
	a.require.NotNil(a.contact, "unexpected nil response")
	a.require.NoError(err, "unexpected error on contact creation")
	return a
}

func (a *assertion) ThatIs(expectedError error) ContactFailureAssertion {
	a.require.ErrorIs(a.err, expectedError)
	return a
}

func (a *assertion) Created(contact *engine.Contact) ContactErrorAssertion {
	a.contact = contact
	return a
}

func (a *assertion) AsNotConfirmed() ContactSuccessAssertion {
	return a.WithStatus(engine.ContactNotConfirmed)
}

func (a *assertion) WithStatus(status engine.ContactStatus) ContactSuccessAssertion {
	a.assert.Equal(status, a.contact.Status)
	return a
}

func (a *assertion) WithFullName(fullName string) ContactSuccessAssertion {
	a.assert.Equal(fullName, a.contact.FullName)
	return a
}

func (a *assertion) WithPubKey(pki string) ContactSuccessAssertion {
	a.assert.Equal(pki, a.contact.PubKey, "counterparty pki invalid")
	return a
}

func (a *assertion) ForUser(user fixtures.User) ContactSuccessAssertion {
	a.assert.Equal(user.XPubID(), a.contact.OwnerXpubID, "contact owner xpub id invalid")
	return a
}

func (a *assertion) ToCounterparty(user fixtures.User) ContactSuccessAssertion {
	a.assert.Equal(user.DefaultPaymail(), a.contact.Paymail, "counterparty paymail invalid")
	return a
}
