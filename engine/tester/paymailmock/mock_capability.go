package paymailmock

import (
	"net/http"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/jarcoal/httpmock"
)

// CapabilityMock is a structure that helps with mocking for a paymail capability.
type CapabilityMock struct {
	name     string
	value    func(name paymailDomainName) any
	endpoint func(name paymailDomainName, c *CapabilityMock) (method string, urlMatcher string, responder httpmock.Responder)
	response httpmock.Responder
}

// WithNotFound will make the capability return the response 404 not found.
func (c *CapabilityMock) WithNotFound() {
	var err error
	c.response, err = httpmock.NewJsonResponder(http.StatusNotFound, obj{"error": "not found"})
	if err != nil {
		panic(spverrors.Wrapf(err, "error creating mocked http response for capability %s", c.name))
	}
}

func (c *CapabilityMock) mockEndpoint(domain paymailDomainName) {
	if c.endpoint != nil {
		httpmock.RegisterResponder(c.endpoint(domain, c))
	}
}

func endpoint(method string, successResponse obj) func(name paymailDomainName, c *CapabilityMock) (method string, urlMatcher string, responder httpmock.Responder) {
	r, err := httpmock.NewJsonResponder(http.StatusOK, successResponse)
	if err != nil {
		panic(err)
	}

	return func(dn paymailDomainName, c *CapabilityMock) (string, string, httpmock.Responder) {
		url, ok := c.value(dn).(string)
		if !ok {
			panic("cannot mock capability without URL in value")
		}
		responder := dynamicResponder(c, r)
		return method, matchingURL(url, string(dn)), responder
	}
}

func dynamicResponder(c *CapabilityMock, defaultResponder httpmock.Responder) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		if c.response != nil {
			return c.response(request)
		}
		return defaultResponder(request)
	}
}

func matchingURL(url string, domain string) string {
	url = strings.ReplaceAll(url, "{alias}", "\\w+")
	url = strings.ReplaceAll(url, "{domain.tld}", domain)
	return "=~" + url
}
