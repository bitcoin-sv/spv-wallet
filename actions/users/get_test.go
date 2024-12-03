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
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		res, _ := client.R().Get("/api/v1/users/current")
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

	t.Run("return xpub info for user (old api)", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		res, _ := client.R().Get("/v1/xpub")
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{.ID}}",
				"created_at": "/.*/",
				"updated_at": "/.*/",
				"current_balance": 0,
				"deleted_at": null,
				"metadata": "*",
				"next_external_num": 1,
				"next_internal_num": 0
			}`, map[string]any{
				"ID": fixtures.Sender.XPubID(),
			})
	})

	t.Run("return xpub info for admin", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		res, _ := client.R().Get("/api/v1/users/current")
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return xpub info for anonymous", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		res, _ := client.R().Get("/api/v1/users/current")
		then.Response(res).IsUnauthorized()
	})
}
