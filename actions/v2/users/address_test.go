package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

func TestCreateAddress(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("create paymail address for user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		res, _ := client.R().
			SetBody(`{
				"paymail": {
					"alias": "test",
					"domain": "spv-wallet.com",
					"publicName": "Test User",
					"avatarURL": "https://example.com/avatar.png"
				}
			}`).
			Post("/api/v2/users/address")

		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"alias": "test",
				"domain": "spv-wallet.com",
				"paymail": "test@spv-wallet.com",
				"publicName": "Test User",
				"avatar": "https://example.com/avatar.png",
				"id": 0
			}`, nil)
	})

	t.Run("return unauthorized for admin", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		res, _ := client.R().
			SetBody(`{
				"paymail": {
					"alias": "test",
					"domain": "spv-wallet.com"
				}
			}`).
			Post("/api/v2/users/address")

		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return unauthorized for anonymous user", func(t *testing.T) {
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		res, _ := client.R().
			SetBody(`{
				"paymail": {
					"alias": "test",
					"domain": "spv-wallet.com"
				}
			}`).
			Post("/api/v2/users/address")

		then.Response(res).IsUnauthorized()
	})
}
