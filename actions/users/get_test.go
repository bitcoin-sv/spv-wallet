package users_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
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
			WithJSONSubsetf(`{
				"id": "%s"
			}`, fixtures.Sender.XPubID())
	})

	t.Run("return xpub info for user (old api)", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		res, _ := client.R().Get("/v1/xpub")
		then.Response(res).
			IsOK().
			WithJSONSubsetf(`{
				"id": "%s"
			}`, fixtures.Sender.XPubID())
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
