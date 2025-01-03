package utxos_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestUserUTXOs(t *testing.T) {
	t.Run("return UTXOs for user", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()

		// and:
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v1/utxos")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`{
			"content": [],
				"page": {
					"number": 1,
					"size": 50,
					"totalElements": 0,
					"totalPages": 0
				}
			}`)

	})

	t.Run("try to return UTXOs for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/utxos")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return UTXOs for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v1/utxos")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
