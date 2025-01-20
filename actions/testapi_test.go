package actions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestUserCurrent(t *testing.T) {
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

}
