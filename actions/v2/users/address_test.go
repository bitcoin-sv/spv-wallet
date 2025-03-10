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

		// when:
		res, _ := client.R().
			SetBody(`{
				"alias": "test",
				"domain": "spv-wallet.com",
				"publicName": "Test User",
				"avatarURL": "https://example.com/avatar.png"
			}`).
			Post("/api/v2/users/address")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"alias": "test",
				"domain": "spv-wallet.com",
				"paymail": "test@spv-wallet.com",
				"publicName": "Test User",
				"avatarURL": "https://example.com/avatar.png"
			}`, nil)
	})

	t.Run("return unauthorized for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(`{
				"alias": "test",
				"domain": "spv-wallet.com"
			}`).
			Post("/api/v2/users/address")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return unauthorized for anonymous user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().
			SetBody(`{
				"alias": "test",
				"domain": "spv-wallet.com"
			}`).
			Post("/api/v2/users/address")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
