package paymailtests

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestPaymailFlow(t *testing.T) {
	//testmode.DevelopmentOnly_SetPostgresModeWithName(t, "spv-test")

	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
		testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	t.Run("Get bsv alias", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		// NOTE: Because testabilities' client has substituted transport,
		// we can make request which looks like it's going to the real server, but instead it goes to the spv-wallet test server.
		// Defining the host in the request is necessary for paymail to work.
		res, _ := client.R().Get("https://example.com/.well-known/bsvalias")

		// then:
		then.Response(res).IsOK().WithJSONf(`{
		  "bsvalias": "1.0",
		  "capabilities": {
			"2a40af698840": "https://example.com/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}",
			"5c55a7fdb7bb": "https://example.com/v1/bsvalias/beef/{alias}@{domain.tld}",
			"5f1323cddf31": "https://example.com/v1/bsvalias/receive-transaction/{alias}@{domain.tld}",
			"6745385c3fc0": false,
			"a9f510c16bde": "https://example.com/v1/bsvalias/verify-pubkey/{alias}@{domain.tld}/{pubkey}",
			"f12f968c92d6": "https://example.com/v1/bsvalias/public-profile/{alias}@{domain.tld}",
			"paymentDestination": "https://example.com/v1/bsvalias/address/{alias}@{domain.tld}",
			"pki": "https://example.com/v1/bsvalias/id/{alias}@{domain.tld}"
		  }
		}`)
	})

	t.Run("Public profile", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// and:
		address := fixtures.Sender.Paymails[0]

		// when:
		res, _ := client.R().Get(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/public-profile/%s",
				address,
			),
		)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"avatar": "{{ matchURL | orEmpty }}",
			"name": "{{ .name }}"
		}`, map[string]any{
			"name": address,
		})
	})

	t.Run("Public profile for not existing paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("https://example.com/v1/bsvalias/public-profile/notexisting@example.com")

		// then:
		then.Response(res).HasStatus(404)
	})

	t.Run("PKI", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// and:
		address := fixtures.Sender.Paymails[0]

		// when:
		res, _ := client.R().Get(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/id/%s",
				address,
			),
		)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"bsvalias": "1.0",
			"handle": "{{ .paymail }}",
			"pubkey": "{{ matchHexWithLength 66 }}"
		}`, map[string]any{
			"paymail": address,
		})
	})
}
