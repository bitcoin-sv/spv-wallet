package transactions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

const (
	transactionsOutlinesRecordURL = "/api/v1/transactions/outlines/record"
	dataOfOpReturnTx              = "hello world"
)

func givenTXWithOpReturn(t *testing.T) fixtures.GivenTXSpec {
	return fixtures.GivenTX(t).
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)
}

func TestOutlinesRecordOpReturn(t *testing.T) {
	t.Run("Record op_return data", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()

		// and:
		client := given.HttpClient().ForUser()

		// and:
		txSpec := givenTXWithOpReturn(t)
		request := `{
			"beef": "` + txSpec.BEEF() + `",
			"annotations": {
				"outputs": {
					"0": {
						"bucket": "data"
					}
				}
			}
		}`

		// and:
		given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: chainmodels.SeenOnNetwork,
		})

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(request).
			Post(transactionsOutlinesRecordURL)

		// then:
		then.Response(res).IsOK()
	})
}

func TestOutlinesRecordOpReturnErrorCases(t *testing.T) {
	givenUnsignedTX := fixtures.GivenTX(t).
		WithoutSigning().
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithP2PKHOutput := fixtures.GivenTX(t).
		WithInput(2).
		WithP2PKHOutput(1)

	tests := map[string]struct {
		request        string
		expectHttpCode int
	}{
		"RecordTransactionOutline for not signed transaction": {
			request: `{
				"beef": "` + givenUnsignedTX.BEEF() + `"
			}`,
			expectHttpCode: 400,
		},
		"RecordTransactionOutline for not a BEEF hex": {
			request: `{
				"beef": "0b3818c665bf28a46""
			}`,
			expectHttpCode: 400,
		},
		"Vout out index as invalid number": {
			request: `{
				"beef": "` + givenTXWithOpReturn(t).BEEF() + `"
				"annotations": {
					"outputs": {
						"invalid-number": {
							"bucket": "data"
						}
					}
				}
			}`,
			expectHttpCode: 400,
		},
		"no-op_return output annotated as data": {
			request: `{
				"beef": "` + givenTxWithP2PKHOutput.BEEF() + `",
				"annotations": {
					"outputs": {
						"0": {
							"bucket": "data"
						}
					}
				}
			}`,
			expectHttpCode: 400,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWallet()
			defer cleanup()

			// and:
			client := given.HttpClient().ForUser()

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				Post(transactionsOutlinesRecordURL)

			// then:
			then.Response(res).HasStatus(test.expectHttpCode)
		})
	}
}
