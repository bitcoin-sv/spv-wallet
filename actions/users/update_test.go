package users_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
)

func TestCurrentUserUpdate(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWallet()
	defer cleanup()

	metadataToUpdate := map[string]any{
		"num": 1234,
		"str": "abc",
	}

	t.Run("update xpub metadata as user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetBody(metadataToUpdate).
			Patch("/api/v1/users/current")

		// then:
		then.Response(res).
			IsOK().
			WithJSONTemplate(`{
				"createdAt": "/.*/",
				"currentBalance": 0,
				"deletedAt": null,
				"id": "{{.ID}}",
				"metadata": {
					"num": 1234,
					"str": "abc"
				},
				"nextExternalNum": 1,
				"nextInternalNum": 0,
				"updatedAt": "/.*/"
			}`, map[string]any{
				"ID": fixtures.Sender.XPubID(),
			})
	})

	t.Run("update xpub metadata for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		// when:
		res, _ := client.R().
			SetBody(metadataToUpdate).
			Patch("/api/v1/users/current")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("update xpub metadata for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		// when:
		res, _ := client.R().
			SetBody(metadataToUpdate).
			Patch("/api/v1/users/current")

		// then
		then.Response(res).IsUnauthorized()
	})
}
