package transactions_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
)

func TestTransactionWithChange(t *testing.T) {
	tests := map[string]struct {
		changeOutputs changeOutputSpecs
	}{
		"one change output": {
			changeOutputs: changeOutputSpecs{
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-1output4d06387d3be7bd26cfe2b5996",
						},
					},
					satoshis: 1000,
				},
			},
		},
		"two change outputs": {
			changeOutputs: changeOutputSpecs{
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-1output4d06387d3be7bd26cfe2b5996",
						},
					},
					satoshis: 600,
				},
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-2output4d06387d3be7bd26cfe2b5996",
						},
					},
					satoshis: 400,
				},
			},
		},
		"one change outputs with longer custom instructions": {
			changeOutputs: changeOutputSpecs{
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-1instruction87d3be7bd26cfe2b5996",
						},
						{
							Type:        "type42",
							Instruction: "1-destination-2instruction7d3be7bd26cfe2b5996",
						},
						{
							Type:        "type42",
							Instruction: "1-destination-3instruction87d3be7bd26cfe2b5996",
						},
					},
					satoshis: 1000,
				},
			},
		},
		"two outputs with the same address": {
			changeOutputs: changeOutputSpecs{
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-1output4d06387d3be7bd26cfe2b5996",
						},
					},
					satoshis: 600,
				},
				{
					customInstructions: bsv.CustomInstructions{
						{
							Type:        "type42",
							Instruction: "1-destination-1output4d06387d3be7bd26cfe2b5996",
						},
					},
					satoshis: 400,
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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
			txSpec := given.Tx().
				WithSender(sender).
				WithInputFromUTXO(sourceTxSpec.TX(), 0).
				WithOPReturn("hello, world")

			for _, changeOutput := range test.changeOutputs {
				txSpec.WithOutputScript(uint64(changeOutput.satoshis), sender.P2PKHLockingScript(changeOutput.customInstructions...))
			}
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
						"outputs": lo.Assign(
							map[string]any{
								"0": map[string]any{
									"bucket": "data",
								},
							},
							test.changeOutputs.toAnnotations(1),
						),
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
		})
	}
}

type changeOutputSpec struct {
	customInstructions bsv.CustomInstructions
	satoshis           bsv.Satoshis
}

func (ch changeOutputSpec) toAnnotation() map[string]any {
	return map[string]any{
		"bucket":             "bsv",
		"customInstructions": ch.customInstructions,
	}
}

type changeOutputSpecs []changeOutputSpec

func (chs changeOutputSpecs) toAnnotations(startingIndex int) map[string]any {
	annotations := make(map[string]any)
	for i, ch := range chs {
		annotations[fmt.Sprintf("%d", i+startingIndex)] = ch.toAnnotation()
	}
	return annotations
}
