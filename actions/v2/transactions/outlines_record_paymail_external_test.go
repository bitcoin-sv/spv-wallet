package transactions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/require"
)

func TestExternalOutgoingTransaction(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	sender := fixtures.Sender
	recipient := fixtures.RecipientExternal

	// and:
	sourceTxSpec := givenForAllTests.Faucet(sender).TopUp(1001)

	// and:
	givenPaymail := givenForAllTests.Paymail()
	externalPaymailHost := givenPaymail.ExternalPaymailHost()
	paymailClientService := givenPaymail.NewPaymailClientService()

	// and:
	externalPaymailHost.WillRespondWithP2PDestinationsWithSats(1000)
	destination, err := paymailClientService.GetP2PDestinations(
		context.Background(),
		&paymail.SanitisedPaymail{
			Alias:   recipient.DefaultPaymail().Alias(),
			Domain:  recipient.DefaultPaymail().Domain(),
			Address: recipient.DefaultPaymail().Address(),
		},
		1000,
	)
	require.NoError(t, err)

	lockingScript, err := script.NewFromHex(destination.Outputs[0].Script)
	require.NoError(t, err)

	// and:
	txSpec := fixtures.GivenTX(t).
		WithSender(sender).
		WithRecipient(recipient).
		WithInputFromUTXO(sourceTxSpec.TX(), 0).
		WithOutputScript(1000, lockingScript)

	// and:
	externalPaymailHost.WillRespondWithP2PWithBEEFCapabilities()

	t.Run("Record new tx outline by sender", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForGivenUser(sender)

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
		}`, txSpec.BEEF(), recipient.DefaultPaymail(), destination.Reference, sender.DefaultPaymail())

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
		then.User(sender).Balance().IsZero()

		// and:
		then.PaymailClient().
			ExternalPaymailHost().
			Called("beef").
			WithRequestJSONMatching(`{
				"beef": "{{ .beef }}",
				"decodedBeef": null,
				"hex": "",
				"metadata": {
					"sender": "{{ .sender }}"
				},
				"reference": "{{ .reference }}"
			}`, map[string]any{
				"beef":      txSpec.BEEF(),
				"sender":    sender.DefaultPaymail(),
				"reference": destination.Reference,
			})
	})

	t.Run("Check sender's operations", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForGivenUser(sender)

		// when:
		res, _ := client.R().Get("/api/v2/operations/search")

		// then:
		then.Response(res).IsOK().WithJSONMatching(`{
			"content": [
				{
					"txID": "{{ .txID }}",
					"createdAt": "{{ matchTimestamp }}",
					"value": {{ .value }},
					"type": "outgoing",
					"counterparty": "{{ .recipient }}",
					"txStatus": "BROADCASTED"
				},
				{{ anything }}
			],
			"page": {
			    "number": 1,
			    "size": 2,
			    "totalElements": 2,
			    "totalPages": 1
			}
		}`, map[string]any{
			"value":     -1001,
			"txID":      txSpec.ID(),
			"recipient": recipient.DefaultPaymail(),
		})
	})
}
