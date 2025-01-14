package operations_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

func TestUserOperations(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	t.Run("return empty operations list for user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v2/operations/search")

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"content": [],
			"page": {
			    "number": 1,
			    "size": 0,
			    "totalElements": 0,
			    "totalPages": 0
			}
		}`, nil)
	})

	t.Run("try return user operations for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/operations/search")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("try return user operations for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v2/operations/search")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
