package paymailserver_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func TestIncomingPaymailRawTX(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
		testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	var testState struct {
		reference     string
		lockingScript *script.Script
	}

	// given:
	given, then := testabilities.NewOf(givenForAllTests, t)
	client := given.HttpClient().ForAnonymous()

	// and:
	recipientPaymail := fixtures.RecipientInternal.DefaultPaymail()
	satoshis := uint64(1000)

	t.Run("step 1 - call p2p-payment-destination", func(t *testing.T) {
		// given:
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
					recipientPaymail,
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
			"satoshis": satoshis,
		})

		// update:
		getter := then.Response(res).JSONValue()
		testState.reference = getter.GetString("reference")

		// and:
		lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
		require.NoError(t, err)
		testState.lockingScript = lockingScript
	})

	t.Run("step 2 - call receive-transaction capability", func(t *testing.T) {
		// TODO: To be implemented in the next PR
	})

	t.Run("step 3 - check balance", func(t *testing.T) {
		// TODO: To be implemented in the next PR
	})
}

func TestIncomingPaymailBeef(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
		testengine.WithNewTransactionFlowEnabled(),
	)
	defer cleanup()

	var testState struct {
		reference     string
		lockingScript *script.Script
		txID          string
	}

	// given:
	given, then := testabilities.NewOf(givenForAllTests, t)
	client := given.HttpClient().ForAnonymous()

	// and:
	recipientPaymail := fixtures.RecipientInternal.DefaultPaymail()
	satoshis := uint64(1000)

	t.Run("step 1 - call p2p-payment-destination", func(t *testing.T) {
		// given:
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
					recipientPaymail,
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
			"satoshis": satoshis,
		})

		// update:
		getter := then.Response(res).JSONValue()
		testState.reference = getter.GetString("reference")

		// and:
		lockingScript, err := script.NewFromHex(getter.GetString("outputs[0]/script"))
		require.NoError(t, err)
		testState.lockingScript = lockingScript
	})

	t.Run("step 2 - call beef capability", func(t *testing.T) {
		// TODO: To be implemented in the next PR
	})

	t.Run("step 3 - check balance", func(t *testing.T) {
		// TODO: To be implemented in the next PR
	})
}
