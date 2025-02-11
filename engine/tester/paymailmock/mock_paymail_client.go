package paymailmock

import (
	"net"
	"net/http"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

// PaymailClientMock is a paymail.Client configured to use mocked endpoints.
type PaymailClientMock struct {
	paymail.ClientInterface
	domains      []paymailDomainName
	capabilities []*CapabilityMock

	mockTransport *httpmock.MockTransport
	sniffer       *httpSniffer
}

// MockClient will return a client for testing purposes
func MockClient(mockTransport *httpmock.MockTransport, domain string, moreDomainNames ...string) *PaymailClientMock {
	domainNames := []paymailDomainName{paymailDomainName(domain)}
	for _, dn := range moreDomainNames {
		domainNames = append(domainNames, paymailDomainName(dn))
	}

	// Create a new client
	newClient, err := paymail.NewClient(
		paymail.WithRequestTracing(),
		paymail.WithDNSTimeout(15*time.Second),
	)

	if err != nil {
		panic(spverrors.Wrapf(err, "error creating mocked paymail client"))
	}

	// Set the HTTP mocking client
	// restyClient -> sniffer -> mockTransport
	sniffer := newHTTPSniffer()
	sniffer.setTransport(mockTransport)

	client := resty.New()
	client.SetTransport(sniffer)

	newClient.WithCustomHTTPClient(client)

	// Build hosts, srv records and ip addresses
	hosts := map[string][]string{}
	records := map[string][]*net.SRV{}
	ipAddresses := map[string][]net.IPAddr{}
	for _, dn := range domainNames {
		name := string(dn)
		hosts[name] = []string{"44.55.66.77", "22.33.44.55", "11.22.33.44"}

		records[paymail.DefaultServiceName+paymail.DefaultProtocol+name] = []*net.SRV{
			{
				Target:   name,
				Port:     paymail.DefaultPort,
				Priority: paymail.DefaultPriority,
				Weight:   paymail.DefaultWeight,
			},
		}

		ipAddresses[name] = []net.IPAddr{
			{IP: net.ParseIP("8.8.8.8"), Zone: "eth0"},
		}
	}

	// Set the custom resolver
	newClient.WithCustomResolver(tester.NewCustomResolver(
		newClient.GetResolver(),
		hosts,
		records,
		ipAddresses,
	))

	return &PaymailClientMock{
		ClientInterface: newClient,
		domains:         domainNames,
		mockTransport:   mockTransport,
		sniffer:         sniffer,
	}
}

// WillRespondWithBasicCapabilities is configuring a client to respond with basic capabilities for all mocked domains.
func (c *PaymailClientMock) WillRespondWithBasicCapabilities() {
	c.mockTransport.Reset()
	c.useBasicCapabilities()
	for _, domain := range c.domains {
		c.exposeCapabilities(domain)
	}
}

// WillRespondWithP2PCapabilities is configuring a client to respond with basic and P2P capabilities for all mocked domains.
func (c *PaymailClientMock) WillRespondWithP2PCapabilities() {
	c.mockTransport.Reset()
	c.useBasicCapabilities()
	c.useP2PCapabilities()
	for _, domain := range c.domains {
		c.exposeCapabilities(domain)
	}
}

// WillRespondWithP2PWithBEEFCapabilities is configuring a client to respond with basic, P2P and BEEF capabilities for all mocked domains.
func (c *PaymailClientMock) WillRespondWithP2PWithBEEFCapabilities() {
	c.mockTransport.Reset()
	c.useBasicCapabilities()
	c.useP2PCapabilities()
	c.useBEEFCapabilities()
	for _, domain := range c.domains {
		c.exposeCapabilities(domain)
	}
}

// WillRespondWithNotFoundOnCapabilities is configuring a client to respond with not found on capabilities for all mocked domains.
func (c *PaymailClientMock) WillRespondWithNotFoundOnCapabilities() {
	c.mockTransport.Reset()
	for _, domain := range c.domains {
		c.exposeErrorOnCapabilities(http.StatusNotFound, domain)
	}
}

// WillRespondWithErrorOnCapabilities is configuring a client to respond with an error on capabilities for all mocked domains.
func (c *PaymailClientMock) WillRespondWithErrorOnCapabilities() {
	c.mockTransport.Reset()
	for _, domain := range c.domains {
		c.exposeErrorOnCapabilities(http.StatusInternalServerError, domain)
	}
}

// WillRespondOnCapability is returning a capability mock for a given capability name.
func (c *PaymailClientMock) WillRespondOnCapability(capabilityName string) *CapabilityMock {
	for _, capability := range c.capabilities {
		if capability.name == capabilityName {
			return capability
		}
	}
	panic(spverrors.Newf("capability %s was't' mocked", capabilityName))
}

// GetMockedServerURL is returning the mocked URL for a given domain if it is mocked.
func (c *PaymailClientMock) GetMockedServerURL(domain string) string {
	return c.findDomain(domain).ServerURL()
}

// GetMockedP2PPaymentDestinationURL is returning the mocked P2P Payment Destination URL for a given domain if it is mocked.
func (c *PaymailClientMock) GetMockedP2PPaymentDestinationURL(domain string) string {
	return c.findDomain(domain).P2PPaymentDestination()
}

// GetMockedP2PTransactionURL is returning the mocked P2P Transaction URL for a given domain if it is mocked.
func (c *PaymailClientMock) GetMockedP2PTransactionURL(domain string) string {
	return c.findDomain(domain).P2PTransaction()
}

// GetMockedBEEFTransactionURL is returning the mocked BEEF Transaction URL for a given domain if it is mocked.
func (c *PaymailClientMock) GetMockedBEEFTransactionURL(domain string) string {
	return c.findDomain(domain).BEEFTransaction()
}

func (c *PaymailClientMock) findDomain(domain string) paymailDomainName {
	for _, dn := range c.domains {
		if string(dn) == domain {
			return dn
		}
	}
	panic(spverrors.Newf("domain %s was't' mocked", domain))
}

// GetCallByRegex is returning the details of a call made to the mocked server by a URL matching a regex.
func (c *PaymailClientMock) GetCallByRegex(r string) *CallDetails {
	return c.sniffer.getCallByRegex(r)
}

func (c *PaymailClientMock) useBasicCapabilities() {
	c.capabilities = append(c.capabilities,
		capabilitySenderValidation(),
		capabilityPki(),
		capabilityPaymentDestination(),
	)
}

func (c *PaymailClientMock) useP2PCapabilities() {
	c.capabilities = append(c.capabilities,
		capabilityP2PTransaction(),
		capabilityP2PPaymentDestination(),
	)
}

func (c *PaymailClientMock) useBEEFCapabilities() {
	c.capabilities = append(c.capabilities,
		capabilityBEEFTransaction(),
	)
}

func (c *PaymailClientMock) exposeCapabilities(domain paymailDomainName) {
	capabilities := make(obj)
	for _, capability := range c.capabilities {
		capabilities[capability.name] = capability.value(domain)
		if capability.endpoint != nil {
			c.mockTransport.RegisterResponder(capability.endpoint(domain, capability))
		}
	}

	bsvAliasResponse := map[string]any{
		paymail.DefaultServiceName: paymail.DefaultBsvAliasVersion,
		"capabilities":             capabilities,
	}

	responder, err := httpmock.NewJsonResponder(http.StatusOK, bsvAliasResponse)
	if err != nil {
		panic(spverrors.Wrapf(err, "error creating mocked capabilities"))
	}
	c.mockTransport.RegisterResponder(http.MethodGet, domain.CapabilitiesURL(), responder)
}

func (c *PaymailClientMock) exposeErrorOnCapabilities(status int, domain paymailDomainName) {
	capabilitiesURL := domain.CapabilitiesURL()
	responder, err := httpmock.NewJsonResponder(status, nil)
	if err != nil {
		panic(spverrors.Wrapf(err, "error creating mocked error capabilities"))
	}
	c.mockTransport.RegisterResponder(http.MethodGet, capabilitiesURL, responder)
}

type obj map[string]any
