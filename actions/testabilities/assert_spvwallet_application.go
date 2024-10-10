package testabilities

import (
	"fmt"
	"mime"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SPVWalletApplicationAssertions interface {
	Response(response *resty.Response) SPVWalletResponseAssertions
}

type SPVWalletResponseAssertions interface {
	IsOK() SPVWalletResponseAssertions
	WithJSONf(expectedFormat string, args ...any)
	IsUnauthorized()
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
	a.assertIsStatus(http.StatusUnauthorized).
		WithJSONf(`{
			"code":"error-unauthorized-auth-header-missing",
			"message":"missing auth header"
		}`)

}

func (a *responseAssertions) IsOK() SPVWalletResponseAssertions {
	return a.assertIsStatus(http.StatusOK)
}

func (a *responseAssertions) assertIsStatus(status int) *responseAssertions {
	a.assert.Equal(status, a.response.StatusCode())
	return a
}

func (a *responseAssertions) WithJSONf(expectedFormat string, args ...any) {
	a.assertJSONContentType()
	a.assertJSONBody(expectedFormat, args...)
}

func (a *responseAssertions) assertJSONContentType() {
	contentType := a.response.Header().Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	a.require.NoError(err, "failed to parse content type")
	a.assert.Equal("application/json", mediaType, "JSON content type expected on response")
}

func (a *responseAssertions) assertJSONBody(expectedFormat string, args ...any) {
	expectedJson := fmt.Sprintf(expectedFormat, args...)
	a.assert.JSONEq(expectedJson, a.response.String())
}
