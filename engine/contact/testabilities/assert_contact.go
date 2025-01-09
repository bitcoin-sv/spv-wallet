package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
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
}

type ContactErrorAssertion interface {
	WithError(err error) ContactFailureAssertion
	WithNoError(err error) ContactSuccessAssertion
}
type ContactAssertion interface {
	Contact(contact *engine.Contact) ContactErrorAssertion
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
	a.require.NoError(err, "record transaction outline has error")
	return a
}

func (a *assertion) ThatIs(expectedError error) ContactFailureAssertion {
	a.require.Nil(a.contact, "unexpected response")
	require.ErrorIs(a.t, a.err, expectedError, "record transaction outline has wrong error")
	return a
}

func (a *assertion) Contact(contact *engine.Contact) ContactErrorAssertion {
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
