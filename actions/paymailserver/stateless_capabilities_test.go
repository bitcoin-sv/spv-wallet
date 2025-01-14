package paymailserver_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestStatelessCapabilities(t *testing.T) {
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
		address := fixtures.RecipientInternal.DefaultPaymail()

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

	t.Run("Get PKI and verify", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// and:
		address := fixtures.RecipientInternal.DefaultPaymail()

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

		// given:
		pki := then.Response(res).JSONValue().GetString("pubkey")

		// when:
		res, _ = client.R().Get(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/verify-pubkey/%s/%s",
				address,
				pki,
			),
		)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"bsvalias": "1.0",
			"handle": "{{ .paymail }}",
			"match": true,
			"pubkey": "{{ .pki }}"
		}`, map[string]any{
			"paymail": address,
			"pki":     pki,
		})
	})

	t.Run("PKI for not existing paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("https://example.com/v1/bsvalias/id/notexisting@example")

		// then:
		then.Response(res).HasStatus(404)
	})

	t.Run("PubKey verify on wrong key", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// and:
		address := fixtures.RecipientInternal.DefaultPaymail()
		wrongPKI := "02561fc133e140526f11438550de3e6cf0ae246a4a5bcd151230652b60124ea1d9"

		// when:
		res, _ := client.R().Get(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/verify-pubkey/%s/%s",
				address,
				wrongPKI,
			),
		)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"bsvalias": "1.0",
			"handle": "{{ .paymail }}",
			"match": false,
			"pubkey": "{{ matchHexWithLength 66 }}"
		}`, map[string]any{
			"paymail": address,
		})
	})

	t.Run("PubKey verify on not existing paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// and:
		notExistingPaymail := "notexisting@example.com"
		wrongPKI := "02561fc133e140526f11438550de3e6cf0ae246a4a5bcd151230652b60124ea1d9"

		// when:
		res, _ := client.R().Get(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/verify-pubkey/%s/%s",
				notExistingPaymail,
				wrongPKI,
			),
		)

		// then:
		then.Response(res).HasStatus(404)
	})
}
