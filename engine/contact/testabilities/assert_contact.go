package testabilities

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/stretchr/testify/require"
)

type assert struct {
	t        testing.TB
	require  *require.Assertions
	response *engine.Contact
}

type ContactFailureAssertion interface {
	WithNilResponse(response *engine.Contact) ContactFailureAssertion
}

type ContactSuccessAssertion interface {
	WithResponse(response *engine.Contact) ContactSuccessAssertion
	WithStatus(status engine.ContactStatus) ContactSuccessAssertion
	WithFullName(fullName string) ContactSuccessAssertion
}

type ContactAssertion interface {
	ErrorIs(err, expectedError error) ContactFailureAssertion
	NoError(err error) ContactSuccessAssertion
}

func then(t testing.TB) ContactAssertion {
	return &assert{
		t:       t,
		require: require.New(t),
	}
}

func (a *assert) NoError(err error) ContactSuccessAssertion {
	a.require.NoError(err, "Record transaction outline has error")
	return a
}

func (a *assert) ErrorIs(err, expectedError error) ContactFailureAssertion {
	require.Error(a.t, err, "Record transaction outline has no error")
	require.ErrorIs(a.t, err, expectedError, "Record transaction outline has wrong error")
	return a
}

func (a *assert) WithResponse(response *engine.Contact) ContactSuccessAssertion {
	a.require.NotNil(response, "unexpected nil response")
	a.response = response
	return a
}

func (a *assert) WithNilResponse(response *engine.Contact) ContactFailureAssertion {
	a.require.Nil(response, "unexpected response")
	return a
}

func (a *assert) WithStatus(status engine.ContactStatus) ContactSuccessAssertion {
	a.require.Equal(status, a.response.Status, fmt.Sprintf("expected status: %s, actual: %s", status, a.response.Status))
	return a
}

func (a *assert) WithFullName(fullName string) ContactSuccessAssertion {
	a.require.Equal(fullName, a.response.FullName, fmt.Sprintf("expected fullName: %s, actual: %s", fullName, a.response.FullName))
	return a
}
