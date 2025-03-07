package transactions_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

const (
	transactionsOutlinesRecordURL = "/api/v2/transactions"
	dataOfOpReturnTx              = "hello world"
)

func givenTXWithOpReturn(t *testing.T) txtestability.TransactionSpec {
	return txtestability.Given(t).Tx().
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)
}

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
		client := given.HttpClient().ForUser()

		// and:
		request := `{
			"hex": "` + txSpec.BEEF() + `",
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
		then.Response(res).
			IsCreated().
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

func TestOutlinesRecordOpReturnErrorCases(t *testing.T) {
	given := testabilities.Given(t)

	givenUnsignedTX := given.Tx().
		WithoutSigning().
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithP2PKHOutput := given.Tx().
		WithInput(2).
		WithP2PKHOutput(1)

	tests := map[string]struct {
		request        string
		expectHttpCode int
		expectedErr    string
	}{
		"RecordTransactionOutline for not signed transaction": {
			request: `{
				"hex": "` + givenUnsignedTX.BEEF() + `"
			}`,
			expectHttpCode: 400,
			expectedErr:    apierror.ExpectedJSON("error-transaction-validation", "transaction validation failed"),
		},
		"RecordTransactionOutline for not a BEEF hex": {
			request: `{
				"hex": "0b3818c665bf28a46""
			}`,
			expectHttpCode: 400,
			expectedErr:    apierror.ExpectedJSON("error-bind-body-invalid", "cannot bind request body"),
		},
		"Vout out index as invalid number": {
			request: `{
				"hex": "` + givenTXWithOpReturn(t).BEEF() + `"
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
				"hex": "` + givenTxWithP2PKHOutput.BEEF() + `",
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
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
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
	txSpec := givenTXWithOpReturn(t)
	request := `{
		"hex": "` + txSpec.BEEF() + `",
		"annotations": {
			"outputs": {
				"0": {
					"bucket": "data"
				}
			}
		}
	}`

	mockTxInfo := func(txStatus chainmodels.TXStatus) chainmodels.TXInfo {
		return chainmodels.TXInfo{
			TxID:     txSpec.ID(),
			TXStatus: txStatus,
		}
	}

	tests := map[string]struct {
		httpStatus      int
		fromBroadcaster chainmodels.TXInfo
	}{
		"500": {
			httpStatus: 500,
		},
		"DoubleSpendAttempted": {
			fromBroadcaster: mockTxInfo(chainmodels.DoubleSpendAttempted),
		},
		"Rejected": {
			fromBroadcaster: mockTxInfo(chainmodels.Rejected),
		},
		"SeenInOrphanMempool": {
			fromBroadcaster: mockTxInfo(chainmodels.SeenInOrphanMempool),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
			defer cleanup()

			// and:
			client := given.HttpClient().ForUser()

			// and:
			given.ARC().WillRespondForBroadcast(test.httpStatus, &test.fromBroadcaster)

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(request).
				Post(transactionsOutlinesRecordURL)

			// then:
			then.Response(res).HasStatus(500).WithJSONf(apierror.ExpectedJSON("error-tx-broadcast", "failed to broadcast transaction"))
		})
	}
}

func TestOutlinesRecordForDifferentTxStatuses(t *testing.T) {
	// given:
	txSpec := givenTXWithOpReturn(t)
	request := `{
			"hex": "` + txSpec.BEEF() + `",
			"annotations": {
				"outputs": {
					"0": {
						"bucket": "data"
					}
				}
			}
		}`

	tests := map[string]struct {
		fromBroadcaster chainmodels.TXStatus
		expectedStatus  string
	}{
		"Mined": {
			fromBroadcaster: chainmodels.Mined,
			expectedStatus:  "MINED",
		},
		"SeenOnNetwork": {
			fromBroadcaster: chainmodels.SeenOnNetwork,
			expectedStatus:  "BROADCASTED",
		},
		"AcceptedByNetwork": {
			fromBroadcaster: chainmodels.AcceptedByNetwork,
			expectedStatus:  "BROADCASTED",
		},
		"SentToNetwork": {
			fromBroadcaster: chainmodels.SentToNetwork,
			expectedStatus:  "BROADCASTED",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
			defer cleanup()

			// and:
			client := given.HttpClient().ForUser()

			// and:
			given.ARC().WillRespondForBroadcast(200, &chainmodels.TXInfo{
				TxID:     txSpec.ID(),
				TXStatus: test.fromBroadcaster,
			})

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(request).
				Post(transactionsOutlinesRecordURL)

			// then:
			then.Response(res).IsCreated()

			// when:
			res, _ = client.R().Get("/api/v2/operations/search")

			// then:
			then.Response(res).IsOK().WithJSONMatching(`{
			"content": [
				{
					"txID": "{{ .txID }}",
					"createdAt": "{{ matchTimestamp }}",
					"value": 0,
					"type": "data",
					"counterparty": "",
					"txStatus": "{{ .expectedStatus }}"
				}
			],
			"page": {
			    "number": 1,
			    "size": 1,
			    "totalElements": 1,
			    "totalPages": 1
				}
		}`, map[string]any{
				"txID":           txSpec.ID(),
				"expectedStatus": test.expectedStatus,
			})
		})
	}
}

func TestRecordOpReturnTwiceByTheSameUser(t *testing.T) {
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

	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("Record op_return data - call %d", i), func(t *testing.T) {
			// given:
			given, then := testabilities.NewOf(givenForAllTests, t)

			// and:
			client := given.HttpClient().ForUser()

			// and:
			request := `{
				"hex": "` + txSpec.BEEF() + `",
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
			then.Response(res).
				IsCreated().
				WithJSONMatching(`{
				"txID": "{{ .txID }}"
			}`, map[string]any{
					"txID": txSpec.ID(),
				})
		})
	}

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
