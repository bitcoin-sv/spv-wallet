package paymailserver_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestIncomingPaymailTX(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
		//testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	var testState struct {
		reference string
	}

	// given:
	given, then := testabilities.NewOf(givenForAllTests, t)
	client := given.HttpClient().ForAnonymous()

	// and:
	address := fixtures.Sender.Paymails[0]

	t.Run("call p2p-payment-destination", func(t *testing.T) {
		// given:
		var requestBody struct {
			Satoshis int `json:"satoshis"`
		}
		requestBody.Satoshis = 1000

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			Post(
				fmt.Sprintf(
					"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
					address,
				),
			)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"outputs": [
			  {
			    "address": "{{ matchAddress }}",
			    "satoshis": {{ .satoshis }},
			    "script": "{{ matchHex }}"
			  }
			],
			"reference": "{{ matchHexWithLength 32 }}"
		}`, map[string]any{
			"satoshis": requestBody.Satoshis,
		})
		// TODO: Question: Do we want to split satoshis into multiple outputs?

		// update:
		testState.reference = then.Response(res).JSONValue().GetString("reference")
	})

	t.Run("call receive-transaction", func(t *testing.T) {

	})
}
