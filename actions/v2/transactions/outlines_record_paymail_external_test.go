package transactions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestExternalOutgoingTransaction(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	sender := fixtures.Sender
	recipient := fixtures.RecipientExternal

	// and:
	sourceTxSpec := given.Faucet(sender).TopUp(1001)

	// and:
	givenPaymail := given.Paymail()
	externalPaymailHost := givenPaymail.ExternalPaymailHost()

	// and:
	lockingScript := recipient.P2PKHLockingScript()
	reference := "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"

	// and:
	txSpec := fixtures.GivenTX(t).
		WithSender(sender).
		WithRecipient(recipient).
		WithInputFromUTXO(sourceTxSpec.TX(), 0).
		WithOutputScript(1000, lockingScript)

	// and:
	externalPaymailHost.WillRespondWithP2PWithBEEFCapabilities()

	// and:
	client := given.HttpClient().ForGivenUser(sender)

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
							"reference": reference,
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

	// and:
	then.User(sender).Operations().Last().
		WithTxID(txSpec.ID()).
		WithTxStatus("BROADCASTED").
		WithValue(-1001).
		WithType("outgoing").
		WithCounterparty(recipient.DefaultPaymail().Address())

	// and:
	then.ExternalPaymailHost().
		ReceivedBeefTransaction(sender.DefaultPaymail().Address(), txSpec.BEEF(), reference)
}
