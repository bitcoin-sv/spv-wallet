package transactions_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
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
	satoshis := bsv.Satoshis(1000)
	sender := fixtures.Sender
	recipient := fixtures.RecipientInternal

	// and:
	sourceTxSpec := givenForAllTests.Faucet(sender).TopUp(satoshis + 1)

	t.Run("During outline preparation - call recipient destination", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAnonymous()

		// and:
		requestBody := map[string]any{
			"satoshis": satoshis,
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
			WithOutputScript(uint64(satoshis), testState.lockingScript)

		// and:
		outline := fmt.Sprintf(`{
          "hex": "%s",
          "format": "BEEF",
          "annotations": {
			"outputs": {
			  "0": {
				"bucket": "bsv",
				"paymail": {
				  "receiver": "%s",
				  "reference": "%s",
				  "sender": "%s"
				}
			  }
			}
		  }
		}`, txSpec.BEEF(), recipient.DefaultPaymail(), testState.reference, sender.DefaultPaymail())

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(outline).
			Post(transactionsOutlinesRecordURL)

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"txID": "{{ .txID }}"
			}`, map[string]any{
				"txID": txSpec.ID(),
			})

		// and:
		thenEng := then.Engine(given.Engine())
		thenEng.User(sender).Balance().IsZero()
		thenEng.User(recipient).Balance().IsEqualTo(satoshis)
	})
}
