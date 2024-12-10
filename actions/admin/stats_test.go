package admin_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
)

func TestGETAdminStats(t *testing.T) {
	t.Run("return unauthorized if not authenticated as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForUser()
		// when:
		res, _ := client.R().Get("/api/v1/admin/stats")

		// then:
		then.Response(res).IsUnauthorizedForUser()
	})

	t.Run("return unauthorized if not authenticated at all", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForAnonymous()
		// when:
		res, _ := client.R().Get("/api/v1/admin/stats")

		// then:
		then.Response(res).IsUnauthorized()
	})

	t.Run("return stats if authenticated as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForAdmin()
		// when:
		res, _ := client.R().Get("/api/v1/admin/stats")

		// then:
		then.Response(res).IsOK().
			WithJSONf(`{
			  "balance": 0,
			  "destinations": 0,
			  "paymailAddresses": 4,
			  "transactions": 0,
			  "transactionsPerDay": null,
			  "utxos": 0,
			  "utxosPerType": null,
			  "xpubs": 4
			}`)
	})
}
