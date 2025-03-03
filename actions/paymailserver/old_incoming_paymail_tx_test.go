package paymailserver_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func TestOldIncomingPaymailRawTX(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
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
	senderPaymail := fixtures.SenderExternal.DefaultPaymail()
	recipientPaymail := fixtures.RecipientInternal.DefaultPaymail()
	satoshis := uint64(1000)
	note := "test note"

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
		// given:
		txSpec := given.Tx().
			WithInput(satoshis+1).
			WithOutputScript(satoshis, testState.lockingScript)

		// and:
		requestBody := map[string]any{
			"hex":       txSpec.RawTX(),
			"reference": testState.reference,
			"metadata": map[string]any{
				"note":   note,
				"sender": senderPaymail,
			},
		}

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			Post(
				fmt.Sprintf(
					"https://example.com/v1/bsvalias/receive-transaction/%s",
					recipientPaymail,
				),
			)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"txid": "{{ .txid }}",
			"note": "{{ .note }}"
		}`, map[string]any{
			"txid": txSpec.ID(),
			"note": note,
		})
	})

	t.Run("step 3 - check balance", func(t *testing.T) {
		// given:
		recipientClient := given.HttpClient().ForGivenUser(fixtures.RecipientInternal)

		// when:
		res, _ := recipientClient.R().Get("/api/v1/users/current")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchID64 }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"currentBalance": {{ .balance }},
				"deletedAt": null,
				"metadata": "*",
				"nextExternalNum": 1,
				"nextInternalNum": 0
			}`, map[string]any{
				"balance": satoshis,
			})
	})
}

func TestOldIncomingPaymailBeef(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
	)
	defer func() {
		// Workaround for the issue with the wallet shutdown
		// There is a `go saveBEEFTxInputs` goroutine that is not finished when the test ends
		time.Sleep(1 * time.Second)
		cleanup()
	}()

	var testState struct {
		reference     string
		lockingScript *script.Script
		txID          string
	}

	// given:
	given, then := testabilities.NewOf(givenForAllTests, t)
	client := given.HttpClient().ForAnonymous()

	// and:
	senderPaymail := fixtures.SenderExternal.DefaultPaymail()
	recipientPaymail := fixtures.RecipientInternal.DefaultPaymail()
	satoshis := uint64(1000)
	note := "test note"

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
		// given:
		txSpec := given.Tx().
			WithInput(satoshis+1).
			WithOutputScript(satoshis, testState.lockingScript)

		// and:
		requestBody := map[string]any{
			"beef":      txSpec.BEEF(),
			"reference": testState.reference,
			"metadata": map[string]any{
				"note":   note,
				"sender": senderPaymail,
			},
		}

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// and;
		given.BHS().WillRespondForMerkleRootsVerify(200, &chainmodels.MerkleRootsConfirmations{
			ConfirmationState: chainmodels.MRConfirmed,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			Post(
				fmt.Sprintf(
					"https://example.com/v1/bsvalias/beef/%s",
					recipientPaymail,
				),
			)

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"txid": "{{ .txid }}",
			"note": "{{ .note }}"
		}`, map[string]any{
			"txid": txSpec.ID(),
			"note": note,
		})
	})

	t.Run("step 3 - check balance", func(t *testing.T) {
		// given:
		recipientClient := given.HttpClient().ForGivenUser(fixtures.RecipientInternal)

		// when:
		res, _ := recipientClient.R().Get("/api/v1/users/current")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchID64 }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"currentBalance": {{ .balance }},
				"deletedAt": null,
				"metadata": "*",
				"nextExternalNum": 1,
				"nextInternalNum": 0
			}`, map[string]any{
				"balance": satoshis,
			})
	})
}

func TestOldAddressResolution(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithDomainValidationDisabled(),
	)
	defer cleanup()

	// given:
	given, then := testabilities.NewOf(givenForAllTests, t)
	client := given.HttpClient().ForAnonymous()

	// and:
	senderPaymail := fixtures.SenderExternal.DefaultPaymail()
	recipientPaymail := fixtures.RecipientInternal.DefaultPaymail()
	satoshis := uint64(1000)

	// and:
	requestBody := map[string]any{
		"dt":           time.Now().UTC().Format(time.RFC3339),
		"senderHandle": senderPaymail,
		"senderName":   "External Sender",
		"purpose":      "P2P",
		"amount":       satoshis,
	}

	// when:
	res, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/address/%s",
				recipientPaymail,
			),
		)

	// then:
	then.Response(res).
		IsOK().
		WithJSONMatching(`{
			"address": "{{ matchAddress }}",
			"output": "{{ matchHex }}"
		}`, nil)
}
