package transactions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

const transactionsOutlinesURL = "/api/v2/transactions/outlines"

func TestPOSTTransactionOutlinesBEEF(t *testing.T) {
	successTestCases := map[string]struct {
		request          string
		responseTemplate string
		responseParams   map[string]any
	}{
		"create transaction outline for op_return strings as default data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for op_return strings data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for op_return hex data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": [ "736f6d65", "20", "64617461" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail without sender": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail()),
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail with sender": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000,
				  "from": "%s"
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail and data": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000,
				  "from": "%s"
				},
				{
				  "type": "op_return",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
			responseTemplate: `{
			  "hex": "{{ matchBEEF }}",
			  "format": "BEEF",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  },
				  "1": {
					"bucket": "data"
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
	}
	for name, test := range successTestCases {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
			defer cleanup()

			// and:
			given.Faucet(fixtures.Sender).TopUp(1_000_000)

			// and:
			client := given.HttpClient().ForUser()

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				Post(transactionsOutlinesURL)

			// then:
			thenResponse := then.Response(res)

			thenResponse.IsOK().
				WithJSONMatching(test.responseTemplate, given.OutlineResponseContext(test.responseParams))

			thenResponse.ContainsValidBEEFHexInField("hex")
		})

		t.Run(name+" with explicit format query", func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
			defer cleanup()

			// and:
			given.Faucet(fixtures.Sender).TopUp(1_000_000)

			// and:
			client := given.HttpClient().ForUser()

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				SetQueryParam("format", "beef").
				Post(transactionsOutlinesURL)

			// then:
			thenResponse := then.Response(res)

			thenResponse.IsOK().
				WithJSONMatching(test.responseTemplate, given.OutlineResponseContext(test.responseParams))

			thenResponse.ContainsValidBEEFHexInField("hex")
		})
	}
}

func TestPOSTTransactionOutlinesRAW(t *testing.T) {
	successTestCases := map[string]struct {
		request          string
		responseTemplate string
		responseParams   map[string]any
	}{
		"create transaction outline for op_return strings as default data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for op_return strings data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for op_return hex data": {
			request: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": [ "736f6d65", "20", "64617461" ]
				}
			  ]
			}`,
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
					}
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail without sender": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail()),
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail with sender": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000,
				  "from": "%s"
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
		"create transaction outline for paymail and data": {
			request: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "%s",
				  "satoshis": 1000,
				  "from": "%s"
				},
				{
				  "type": "op_return",
				  "data": [ "some", " ", "data" ]
				}
			  ]
			}`, fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
			responseTemplate: `{
			  "hex": "{{ matchHex }}",
			  "format": "RAW",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "{{ .ReceiverPaymail }}",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "{{ .SenderPaymail }}"
					}
				  },
				  "1": {
					"bucket": "data"
				  }
				},
				"inputs": {
				  "0": {
				    "customInstructions": {{ .CustomInstructions }}
				  }
				}
			  }
			}`,
		},
	}
	for name, test := range successTestCases {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
			defer cleanup()

			// and:
			given.Faucet(fixtures.Sender).TopUp(1_000_000)

			// and:
			client := given.HttpClient().ForUser()

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.request).
				SetQueryParam("format", "raw").
				Post(transactionsOutlinesURL)

			// then:
			thenResponse := then.Response(res)

			thenResponse.IsOK().
				WithJSONMatching(test.responseTemplate, given.OutlineResponseContext(test.responseParams))

			thenResponse.ContainsValidRawTxHexInField("hex").
				WithOutValues(0, 999_999)
		})
	}
}

func TestPOSTTransactionOutlinesErrors(t *testing.T) {
	t.Run("not allowed for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
		defer cleanup()

		// and:
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Post(transactionsOutlinesURL)

		// then:
		then.Response(res).IsUnauthorized()
	})

	t.Run("not allowed for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
		defer cleanup()

		// and:
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Post(transactionsOutlinesURL)

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("Bad Request: when user has no paymail address", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
		defer cleanup()

		// and:
		client := given.HttpClient().ForGivenUser(fixtures.UserWithoutPaymail)

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "recipient@example.com",
				  "satoshis": 1
				}
			  ]
			}`).
			Post(transactionsOutlinesURL)

		// then:
		then.Response(res).IsBadRequest().WithJSONf(apierror.ExpectedJSON("error-tx-spec-paymail-address-no-default", "cannot choose paymail address of the sender"))
	})

	t.Run("Bad Request: no body", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
		defer cleanup()

		// and:
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			Post(transactionsOutlinesURL)

		// then:
		then.Response(res).IsBadRequest().WithJSONf(apierror.CannotBindBodyJSON)
	})

	badRequestTestCases := map[string]struct {
		json           string
		expectedStatus int
		expectedErr    string
	}{
		"Bad Request: Empty request": {
			json:           `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("tx-spec-output-required", "transaction outline requires at least one output"),
		},
		"Bad Request: Empty outputs": {
			json: `{
			  "outputs": []
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("tx-spec-output-required", "transaction outline requires at least one output"),
		},
		"Bad Request: Unsupported output type": {
			json: `{
			  "outputs": [
				{
				  "type": "unsupported"
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: OP_RETURN output without data": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: OP_RETURN output with empty data list": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "data": []
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
		},
		"Bad Request: OP_RETURN output with unknown data type": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "unknown",
				  "data": [ "hello world" ]
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: OP_RETURN strings output empty data list": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": []
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
		},
		"Bad Request: OP_RETURN strings output with string instead of array as data": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": "hello world"
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: OP_RETURN hexes output with empty data list": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": []
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
		},
		"Bad Request: OP_RETURN hexes output with invalid hex": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": ["invalid hex"]
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("failed-to-decode-hex", "failed to decode hex"),
		},
		"Bad Request: OP_RETURN hexes output with single hex instead of list": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": "0"
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: Paymail output with negative satoshis": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver@example.com",
				  "satoshis": -1
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.CannotBindBodyJSON,
		},
		"Bad Request: Paymail output without paymail address": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "satoshis": 1
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
		},
		"Bad Request: Paymail output with only alias without domain": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver",
				  "satoshis": 1
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
		},
		"Bad Request: Paymail output with only domain without alias": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "@example.com",
				  "satoshis": 1
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
		},
		"Bad Request: Paymail output with sender address not existing in our system": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "recipient@example.com",
				  "from": "not_existing_alias@example.com",
				  "satoshis": 1
				}
			  ]
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("error-paymail-address-invalid-sender", "sender paymail address is invalid"),
		},
		"Bad Request: Paymail output with sender address not belongin to that user": {
			json: fmt.Sprintf(`{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "recipient@example.com",
				  "from": "%s",
				  "satoshis": 1
				}
			  ]
			}`, fixtures.UserWithMorePaymails.DefaultPaymail()),
			expectedStatus: http.StatusBadRequest,
			expectedErr:    apierror.ExpectedJSON("error-paymail-address-invalid-sender", "sender paymail address is invalid"),
		},
		"Unprocessable: User has not enough funds": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "data": [ "1" ]
				}
			  ]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedErr:    apierror.ExpectedJSON("tx-outline-not-enough-funds", "not enough funds to make the transaction"),
		},
	}
	for name, test := range badRequestTestCases {
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
				SetBody(test.json).
				Post(transactionsOutlinesURL)

			// then:
			then.Response(res).
				HasStatus(test.expectedStatus).
				WithJSONf(test.expectedErr)
		})
	}
}
