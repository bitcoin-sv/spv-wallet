package actions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestOApi(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	t.Run("For user", func(t *testing.T) {
		// given:
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v2/testapi?something=1")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"xpubID": "{{ .id }}",
				"something": "1"
			}`, map[string]any{"id": fixtures.Sender.XPubID()})
	})

	t.Run("For admin", func(t *testing.T) {
		// given:
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/admin/testapi")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"xpubID": "admin",
				"something": "something"
			}`, nil)

	})

	t.Run("Transactions - testing polymorphism", func(t *testing.T) {
		// given:
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"outputs": []map[string]any{
					{
						"dataType": "op_return",
						"data":     []string{"1", "2", "3"},
					},
					{
						"dataType": "paymail",
						"to":       "user@exapmle.com",
						"satoshis": 1,
					},
				},
			}).
			Post("/api/v2/testpolymorphism")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"xpubID": "{{ .id }}",
				"something": "1, 2, 3, user@exapmle.com"
			}`, map[string]any{"id": fixtures.Sender.XPubID()})
	})

}
