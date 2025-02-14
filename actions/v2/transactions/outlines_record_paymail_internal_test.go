package transactions_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func TestInternalOutgoingTransaction(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	var testState struct {
		reference     string
		lockingScript *script.Script
	}

	// and:
	sender := fixtures.Sender
	recipient := fixtures.RecipientInternal

	// and:
	sourceTxSpec := givenForAllTests.Faucet(sender).TopUp(1001)

	t.Run("During outline preparation - call recipient destination", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAnonymous()

		// and:
		requestBody := map[string]any{
			"satoshis": 1000,
		}

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			Post(
				fmt.Sprintf(
					"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
					recipient.DefaultPaymail(),
				),
			)

		// then:
		then.Response(res).IsOK()

		// update:
		getter := then.Response(res).JSONValue()
		testState.reference = getter.GetString("reference")

		// and:
		lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
		require.NoError(t, err)
		testState.lockingScript = lockingScript
	})

	t.Run("Record new tx outline by sender", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForGivenUser(sender)

		// and:
		txSpec := fixtures.GivenTX(t).
			WithSender(sender).
			WithRecipient(recipient).
			WithInputFromUTXO(sourceTxSpec.TX(), 0).
			WithOutputScript(1000, testState.lockingScript)

		// and:
		given.ARC().WillRespondForBroadcastWithSeenOnNetwork(txSpec.ID())

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]any{
				"hex":    txSpec.BEEF(),
				"format": "BEEF",
				"annotations": map[string]any{
					"outputs": map[string]any{
						"0": map[string]any{
							"bucket": "bsv",
							"paymail": map[string]any{
								"receiver":  recipient.DefaultPaymail(),
								"reference": testState.reference,
								"sender":    sender.DefaultPaymail(),
							},
						},
					},
				},
			}).
			Post(transactionsOutlinesRecordURL)

		// then:
		then.Response(res).
			IsCreated().
			WithJSONMatching(`{
				"txID": "{{ .txID }}"
			}`, map[string]any{
				"txID": txSpec.ID(),
			})

		// and:
		then.User(sender).Balance().IsZero()
		then.User(sender).Operations().Last().
			WithTxID(txSpec.ID()).
			WithTxStatus("BROADCASTED").
			WithValue(-1001).
			WithType("outgoing").
			WithCounterparty(recipient.DefaultPaymail().Address())

		// and:
		then.User(recipient).Balance().IsEqualTo(1000)
		then.User(recipient).Operations().Last().
			WithTxID(txSpec.ID()).
			WithTxStatus("BROADCASTED").
			WithValue(1000).
			WithType("incoming").
			WithCounterparty(sender.DefaultPaymail().Address())
	})
}
