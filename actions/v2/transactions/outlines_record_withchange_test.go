package transactions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

func TestTransactionWithChange(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	sender := fixtures.Sender

	// and:
	sourceTxSpec := given.Faucet(sender).TopUp(1001)

	// and:
	changeCustomInstr := []bsv.CustomInstruction{
		{
			Type:        "type42",
			Instruction: "1-destination-cb15dc54d06387d3be7bd26cfe2b5996",
		},
	}

	// and:
	txSpec := fixtures.GivenTX(t).
		WithSender(sender).
		WithInputFromUTXO(sourceTxSpec.TX(), 0).
		WithOPReturn("hello, world").
		WithOutputScript(1000, sender.P2PKHLockingScript(changeCustomInstr...))

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
						"bucket": "data",
					},
					"1": map[string]any{
						"bucket":             "bsv",
						"customInstructions": changeCustomInstr,
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
	then.User(sender).Balance().IsEqualTo(1000)

	// and:
	then.User(sender).Operations().Last().
		WithTxID(txSpec.ID()).
		WithTxStatus("BROADCASTED").
		WithValue(-1).
		WithType("outgoing").
		WithNoCounterparty()

}
