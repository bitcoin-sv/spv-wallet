package transactions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestDoubleSpending(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	sourceTxSpec := givenForAllTests.Faucet(fixtures.Sender).TopUp(1000)

	t.Run("Spending the UTXO", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		txSpec := given.Tx().
			WithSender(fixtures.Sender).
			WithInputFromUTXO(sourceTxSpec.TX(), 0).
			WithOPReturn("hello world")

		// and:
		client := given.HttpClient().ForUser()

		// and:
		request := `{
			"hex": "` + txSpec.BEEF() + `",
			"annotations": {
				"outputs": {
					"0": {
						"bucket": "data"
					}
				}
			}
		}`

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(request).
			Post(transactionsOutlinesRecordURL)

		// then:
		then.Response(res).
			IsCreated().
			WithJSONMatching(`{
				"txID": "{{ .txID }}"
			}`, map[string]any{
				"txID": txSpec.ID(),
			})
	})

	t.Run("Double spend attempt", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		txSpec := given.Tx().
			WithSender(fixtures.Sender).
			WithInputFromUTXO(sourceTxSpec.TX(), 0).
			WithOPReturn("other data")

		// and:
		client := given.HttpClient().ForUser()

		// and:
		request := `{
			"hex": "` + txSpec.BEEF() + `",
			"annotations": {
				"outputs": {
					"0": {
						"bucket": "data"
					}
				}
			}
		}`

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(request).
			Post(transactionsOutlinesRecordURL)

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-utxo-spent", "UTXO is already spent"))
	})
}
