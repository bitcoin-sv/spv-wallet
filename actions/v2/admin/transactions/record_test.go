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
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	recordTransactionOutlineForUserURL = "/api/v2/admin/transactions/record"
	dataOfOpReturnTx                   = "hello world"
)

func TestOutlinesRecordOpReturn(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	ownedTransaction := givenForAllTests.Faucet(fixtures.Sender).TopUp(1000)

	// and:
	txSpec := givenForAllTests.Tx().
		WithSender(fixtures.Sender).
		WithInputFromUTXO(ownedTransaction.TX(), 0).
		WithOPReturn(dataOfOpReturnTx)

	t.Run("Record op_return data", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAdmin()

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]any{
				"hex": txSpec.BEEF(),
				"annotations": map[string]any{
					"outputs": map[string]any{
						"0": map[string]any{
							"bucket": "data",
						},
					},
				},
				"userID": fixtures.Sender.ID(),
			}).
			Post(recordTransactionOutlineForUserURL)

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"txID": "{{ .txID }}"
			}`, map[string]any{
				"txID": txSpec.ID(),
			})
	})

	t.Run("get operations", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForUser()

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
					"counterparty": "{{ .sender }}",
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
			"value":  -1000,
			"txID":   txSpec.ID(),
			"sender": "",
		})
	})

	t.Run("Get stored data", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		outpoint := bsv.Outpoint{TxID: txSpec.ID(), Vout: 0}

		// and:
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Get("/api/v2/data/" + outpoint.String())

		// then:
		then.Response(res).
			IsOK().WithJSONMatching(`{
				"id": "{{ .outpoint }}",
				"blob": "{{ .blob }}"
			}`, map[string]any{
			"outpoint": outpoint.String(),
			"blob":     dataOfOpReturnTx,
		})
	})
}

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
	sourceTxSpec := givenForAllTests.Faucet(sender).TopUp(1201)

	t.Run("During outline preparation - call recipient destination", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAnonymous()

		// and:
		requestBody := map[string]any{
			"satoshis": 200,
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

	t.Run("Record new tx outline with change for sender by admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAdmin()

		// and:
		txSpec := given.Tx().
			WithSender(sender).
			WithRecipient(recipient).
			WithInputFromUTXO(sourceTxSpec.TX(), 0).
			WithOutputScript(200, testState.lockingScript)

		// and:
		changeOutputs := changeOutputSpecs{
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
				satoshis: 200,
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
			{
				customInstructions: bsv.CustomInstructions{
					{
						Type:        "type42",
						Instruction: "1-destination-1output4d06387d3be7bd26cfe2b5996",
					},
				},
				satoshis: 400,
			},
		}
		for _, changeOutput := range changeOutputs {
			txSpec.WithOutputScript(uint64(changeOutput.satoshis), sender.P2PKHLockingScript(changeOutput.customInstructions...))
		}

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
								"bucket": "bsv",
								"paymail": map[string]any{
									"receiver":  recipient.DefaultPaymail(),
									"reference": testState.reference,
									"sender":    sender.DefaultPaymail(),
								},
							},
						},
						changeOutputs.toAnnotations(1),
					),
				},
				"userID": sender.ID(),
			}).
			Post(recordTransactionOutlineForUserURL)

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
		then.User(sender).Operations().Last().
			WithTxID(txSpec.ID()).
			WithTxStatus("BROADCASTED").
			WithValue(-201).
			WithType("outgoing").
			WithCounterparty(recipient.DefaultPaymail().Address())

		// and:
		then.User(recipient).Balance().IsEqualTo(200)
		then.User(recipient).Operations().Last().
			WithTxID(txSpec.ID()).
			WithTxStatus("BROADCASTED").
			WithValue(200).
			WithType("incoming").
			WithCounterparty(sender.DefaultPaymail().Address())
	})
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
