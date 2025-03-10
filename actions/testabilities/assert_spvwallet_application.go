package testabilities

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SPVWalletApplicationAssertions interface {
	Response(response *resty.Response) SPVWalletResponseAssertions
	User(user fixtures.User) SPVWalletAppUserAssertions
	ExternalPaymailHost() testpaymail.PaymailExternalAssertions
	ARC() testengine.ARCAssertions
}

type SPVWalletResponseAssertions interface {
	IsOK() SPVWalletResponseAssertions
	IsCreated() SPVWalletResponseAssertions
	HasStatus(status int) SPVWalletResponseAssertions
	WithJSONf(expectedFormat string, args ...any)
	WithJSONMatching(expectedTemplateFormat string, params map[string]any)
	JSONValue() JsonValueGetter
	// IsUnauthorized asserts that the response status code is 401 and the error is about lack of authorization.
	IsUnauthorized()
	// IsUnauthorizedForAdmin asserts that the response status code is 401 and the error is that admin is not authorized to use the endpoint.
	IsUnauthorizedForAdmin()
	// IsUnauthorizedForUser asserts that the response status code is 401 and the error is that only admin is authorized to use this endpoint.
	IsUnauthorizedForUser()
	IsBadRequest() SPVWalletResponseAssertions
}

type JsonValueGetter interface {
	GetString(xpath string) string
	GetAsType(xpath string, target any)
	GetField(xpath string) any
	GetInt(xpath string) int
}

func Then(t testing.TB, app SPVWalletApplicationFixture) SPVWalletApplicationAssertions {
	return &appAssertions{
		t:                t,
		appFixtures:      app,
		engineAssertions: testengine.Then(t, app.EngineFixture()),
	}
}

type appAssertions struct {
	t                testing.TB
	engineAssertions testengine.EngineAssertions
	appFixtures      SPVWalletApplicationFixture
}

func (a *appAssertions) Response(response *resty.Response) SPVWalletResponseAssertions {
	return &responseAssertions{
		t:        a.t,
		require:  require.New(a.t),
		assert:   assert.New(a.t),
		response: response,
	}
}

func (a *appAssertions) User(user fixtures.User) SPVWalletAppUserAssertions {
	return &userAssertions{
		userClient: a.appFixtures.HttpClient().ForGivenUser(user),
		t:          a.t,
		require:    require.New(a.t),
	}
}

func (a *appAssertions) ExternalPaymailHost() testpaymail.PaymailExternalAssertions {
	return a.engineAssertions.ExternalPaymailHost()
}

func (a *appAssertions) ARC() testengine.ARCAssertions {
	return a.engineAssertions.ARC()
}

type responseAssertions struct {
	t        testing.TB
	require  *require.Assertions
	assert   *assert.Assertions
	response *resty.Response
}

func (a *responseAssertions) Response(response *resty.Response) SPVWalletResponseAssertions {
	a.t.Helper()
	a.require.NotNil(response, "unexpected nil response")
	a.response = response
	return a
}

func (a *responseAssertions) IsUnauthorized() {
	a.t.Helper()
	a.HasStatus(http.StatusUnauthorized).
		WithJSONf(apierror.MissingAuthHeaderJSON)

}

func (a *responseAssertions) IsUnauthorizedForAdmin() {
	a.t.Helper()
	a.HasStatus(http.StatusUnauthorized).
		WithJSONf(apierror.AdminNotAuthorizedJSON)
}

func (a *responseAssertions) IsUnauthorizedForUser() {
	a.t.Helper()
	a.HasStatus(http.StatusUnauthorized).
		WithJSONf(apierror.UserNotAuthorizedJSON)
}

func (a *responseAssertions) IsOK() SPVWalletResponseAssertions {
	a.t.Helper()
	return a.HasStatus(http.StatusOK)
}

func (a *responseAssertions) IsCreated() SPVWalletResponseAssertions {
	a.t.Helper()
	return a.HasStatus(http.StatusCreated)
}

func (a *responseAssertions) IsBadRequest() SPVWalletResponseAssertions {
	a.t.Helper()
	return a.HasStatus(http.StatusBadRequest)
}

func (a *responseAssertions) HasStatus(status int) SPVWalletResponseAssertions {
	a.t.Helper()
	a.assert.Equal(status, a.response.StatusCode())
	return a
}

func (a *responseAssertions) WithJSONf(expectedFormat string, args ...any) {
	a.t.Helper()
	a.assertJSONContentType()
	a.assertJSONBody(expectedFormat, args...)
}

func (a *responseAssertions) WithJSONMatching(expectedTemplateFormat string, params map[string]any) {
	a.t.Helper()
	a.assertJSONContentType()
	jsonrequire.Match(a.t, expectedTemplateFormat, params, a.response.String())
}

func (a *responseAssertions) JSONValue() JsonValueGetter {
	a.t.Helper()
	a.assertJSONContentType()
	return jsonrequire.NewGetterWithJSON(a.t, a.response.String())
}

func (a *responseAssertions) assertJSONContentType() {
	a.t.Helper()
	contentType := a.response.Header().Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	a.require.NoError(err, "Cannot validate Content-Type of response because of error")
	a.assert.Equal("application/json", mediaType, "JSON content type expected on response")
}

func (a *responseAssertions) assertJSONBody(expectedFormat string, args ...any) {
	a.t.Helper()
	expectedJson := fmt.Sprintf(expectedFormat, args...)
	assertJSONEq(a.t, expectedJson, a.response.String())
}

// assertJSONEq compares two JSON strings and fails the test if they are not equal.
// It is alternative to assert.JSONEq which provides a diff between structs decoded from json, not JSON strings itself, so it's harder to find a difference.
func assertJSONEq(t testing.TB, expected, actual string) {
	t.Helper()
	var expectedJSONValue, actualJSONValue interface{}
	if err := json.Unmarshal([]byte(expected), &expectedJSONValue); err != nil {
		require.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()))
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONValue); err != nil {
		require.Fail(t, fmt.Sprintf("Input value ('%s') is not valid json.\nJSON parsing error: '%s'", actual, err.Error()))
	}

	if assert.ObjectsAreEqual(expectedJSONValue, actualJSONValue) {
		return
	}

	// We want a message that shows the diff between the two JSON strings not the decoded objects.
	expectedJsonString := expected
	actualJSONString := actual

	// Try to unify the JSON strings to make the diff more readable.
	marshaled, err := json.MarshalIndent(expectedJSONValue, "", "  ")
	if err == nil {
		expectedJsonString = string(marshaled)
	}

	marshaled, err = json.MarshalIndent(actualJSONValue, "", "  ")
	if err == nil {
		actualJSONString = string(marshaled)
	}

	assert.Equal(t, expectedJsonString, actualJSONString)
}
