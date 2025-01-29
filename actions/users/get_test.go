package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestCurrentUserGet(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWallet()
	defer cleanup()

	t.Run("return xpub info for user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v1/users/current")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{.ID}}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"currentBalance": 0,
				"deletedAt": null,
				"metadata": "*",
				"nextExternalNum": 1,
				"nextInternalNum": 0
			}`, map[string]any{
				"ID": fixtures.Sender.XPubID(),
			})
	})

	t.Run("return xpub info for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/users/current")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return xpub info for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v1/users/current")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
