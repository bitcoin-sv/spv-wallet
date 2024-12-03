package testabilities

import (
	"encoding/json"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mime"
	"net/http"
	"testing"
)

type SPVWalletApplicationAssertions interface {
	Response(response *resty.Response) SPVWalletResponseAssertions
}

type SPVWalletResponseAssertions interface {
	IsOK() SPVWalletResponseAssertions
	HasStatus(status int) SPVWalletResponseAssertions
	WithJSONf(expectedFormat string, args ...any)
	WithJSONMatching(expectedTemplateFormat string, params map[string]any)
	// IsUnauthorized asserts that the response status code is 401 and the error is about lack of authorization.
	IsUnauthorized()
	// IsUnauthorizedForAdmin asserts that the response status code is 401 and the error is that admin is not authorized to use the endpoint.
	IsUnauthorizedForAdmin()
	IsBadRequest() SPVWalletResponseAssertions
}

func Then(t testing.TB) SPVWalletApplicationAssertions {
	return &responseAssertions{
		t:       t,
		require: require.New(t),
		assert:  assert.New(t),
	}
}

type responseAssertions struct {
	t        testing.TB
	require  *require.Assertions
	assert   *assert.Assertions
	response *resty.Response
}

func (a *responseAssertions) Response(response *resty.Response) SPVWalletResponseAssertions {
	a.require.NotNil(response, "unexpected nil response")
	a.response = response
	return a
}

func (a *responseAssertions) IsUnauthorized() {
	a.HasStatus(http.StatusUnauthorized).
		WithJSONf(apierror.MissingAuthHeaderJSON)

}

func (a *responseAssertions) IsUnauthorizedForAdmin() {
	a.HasStatus(http.StatusUnauthorized).
		WithJSONf(apierror.AdminNotAuthorizedJSON)
}

func (a *responseAssertions) IsOK() SPVWalletResponseAssertions {
	return a.HasStatus(http.StatusOK)
}

func (a *responseAssertions) IsBadRequest() SPVWalletResponseAssertions {
	return a.HasStatus(http.StatusBadRequest)
}

func (a *responseAssertions) HasStatus(status int) SPVWalletResponseAssertions {
	a.assert.Equal(status, a.response.StatusCode())
	return a
}

func (a *responseAssertions) WithJSONf(expectedFormat string, args ...any) {
	a.assertJSONContentType()
	a.assertJSONBody(expectedFormat, args...)
}

func (a *responseAssertions) WithJSONMatching(expectedTemplateFormat string, params map[string]any) {
	a.assertJSONContentType()
	jsonrequire.Match(a.t, expectedTemplateFormat, params, a.response.String())
}

func (a *responseAssertions) assertJSONContentType() {
	contentType := a.response.Header().Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	a.require.NoError(err, "Cannot validate Content-Type of response because of error")
	a.assert.Equal("application/json", mediaType, "JSON content type expected on response")
}

func (a *responseAssertions) assertJSONBody(expectedFormat string, args ...any) {
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
