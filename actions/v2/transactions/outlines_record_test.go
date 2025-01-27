package transactions_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

const (
	transactionsOutlinesRecordURL = "/api/v2/transactions"
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
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithNewTransactionFlowEnabled())
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
		expectedErr    string
	}{
		"RecordTransactionOutline for not signed transaction": {
			request: `{
				"beef": "` + givenUnsignedTX.BEEF() + `"
			}`,
			expectHttpCode: 400,
			expectedErr:    apierror.ExpectedJSON("error-transaction-validation", "transaction validation failed"),
		},
		"RecordTransactionOutline for not a BEEF hex": {
			request: `{
				"beef": "0b3818c665bf28a46""
			}`,
			expectHttpCode: 400,
			expectedErr:    apierror.ExpectedJSON("error-bind-body-invalid", "cannot bind request body"),
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
			expectedErr:    apierror.ExpectedJSON("error-bind-body-invalid", "cannot bind request body"),
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
			expectedErr:    apierror.ExpectedJSON("error-annotation-mismatch", "annotation mismatch"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithNewTransactionFlowEnabled())
			defer cleanup()

			// and:
			client := given.HttpClient().ForUser()

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				Post(transactionsOutlinesRecordURL)

			// then:
			then.Response(res).HasStatus(test.expectHttpCode).WithJSONf(test.expectedErr)
		})
	}
}

func TestOutlinesRecordOpReturnOnBroadcastError(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithNewTransactionFlowEnabled())
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
	given.ARC().WillRespondForBroadcast(500, &chainmodels.TXInfo{})

	// when:
	res, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(transactionsOutlinesRecordURL)

	// then:
	then.Response(res).HasStatus(500).WithJSONf(apierror.ExpectedJSON("error-tx-broadcast", "failed to broadcast transaction"))
}
