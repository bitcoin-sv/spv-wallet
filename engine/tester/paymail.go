package tester

import (
	"net"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/tester"
)

// PaymailMockClient will return a client for testing purposes
func PaymailMockClient(domainNames []string) (paymail.ClientInterface, error) {

	// Create a new client
	newClient, err := paymail.NewClient(
		paymail.WithRequestTracing(),
		paymail.WithDNSTimeout(15*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// Set the HTTP mocking client
	newClient.WithCustomHTTPClient(tester.MockResty())

	// Build hosts, srv records and ip addresses
	hosts := map[string][]string{}
	records := map[string][]*net.SRV{}
	ipAddresses := map[string][]net.IPAddr{}
	for _, name := range domainNames {
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
	return newClient, nil
}
