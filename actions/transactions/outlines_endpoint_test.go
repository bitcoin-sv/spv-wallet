package transactions_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

const transactionsOutlinesURL = "/api/v1/transactions/outlines"

func TestPOSTTransactionOutlines(t *testing.T) {
	successTestCases := map[string]struct {
		request  string
		response string
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
			response: `{
			  "beef": "0100beef000100000000000100000000000000000e006a04736f6d65012004646174610000000000",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
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
			response: `{
			  "beef": "0100beef000100000000000100000000000000000e006a04736f6d65012004646174610000000000",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
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
			response: `{
			  "beef": "0100beef000100000000000100000000000000000e006a04736f6d65012004646174610000000000",
			  "annotations": {
				"outputs": {
					"0": {
						  "bucket": "data"
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
			response: fmt.Sprintf(`{
			  "beef": "0100beef0001000000000001e8030000000000001976a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac0000000000",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "%s",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "%s"
					}
				  }
				}
			  }
			}`,
				fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
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
			response: fmt.Sprintf(`{
			  "beef": "0100beef0001000000000001e8030000000000001976a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac0000000000",
			  "annotations": {
				"outputs": {
				  "0": {
					"bucket": "bsv",
					"paymail": {
					  "receiver": "%s",
					  "reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
					  "sender": "%s"
					}
				  }
				}
			  }
			}`,
				fixtures.RecipientExternal.DefaultPaymail(),
				fixtures.Sender.DefaultPaymail(),
			),
		},
	}
	for name, test := range successTestCases {
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
				Post(transactionsOutlinesURL)

			// then:
			then.Response(res).IsOK().WithJSONf(test.response)
		})
	}

	t.Run("not allowed for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
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
		cleanup := given.StartedSPVWallet()
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
		cleanup := given.StartedSPVWallet()
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
		cleanup := given.StartedSPVWallet()
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
		json        string
		expectedErr string
	}{
		"Bad Request: Empty request": {
			json:        `{}`,
			expectedErr: apierror.ExpectedJSON("tx-spec-output-required", "transaction outline requires at least one output"),
		},
		"Bad Request: Empty outputs": {
			json: `{
			  "outputs": []
			}`,
			expectedErr: apierror.ExpectedJSON("tx-spec-output-required", "transaction outline requires at least one output"),
		},
		"Bad Request: Unsupported output type": {
			json: `{
			  "outputs": [
				{
				  "type": "unsupported"
				}
			  ]
			}`,
			expectedErr: apierror.CannotBindBodyJSON,
		},
		"Bad Request: OP_RETURN output without data": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				}
			  ]
			}`,
			expectedErr: apierror.CannotBindBodyJSON,
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
			expectedErr: apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
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
			expectedErr: apierror.CannotBindBodyJSON,
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
			expectedErr: apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
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
			expectedErr: apierror.CannotBindBodyJSON,
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
			expectedErr: apierror.ExpectedJSON("tx-spec-op-return-data-required", "data is required for OP_RETURN output"),
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
			expectedErr: apierror.ExpectedJSON("failed-to-decode-hex", "failed to decode hex"),
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
			expectedErr: apierror.CannotBindBodyJSON,
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
			expectedErr: apierror.CannotBindBodyJSON,
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
			expectedErr: apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
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
			expectedErr: apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
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
			expectedErr: apierror.ExpectedJSON("error-paymail-address-invalid-receiver", "receiver paymail address is invalid"),
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
			expectedErr: apierror.ExpectedJSON("error-paymail-address-invalid-sender", "sender paymail address is invalid"),
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
			expectedErr: apierror.ExpectedJSON("error-paymail-address-invalid-sender", "sender paymail address is invalid"),
		},
	}
	for name, test := range badRequestTestCases {
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
				SetBody(test.json).
				Post(transactionsOutlinesURL)

			// then:
			then.Response(res).IsBadRequest().WithJSONf(test.expectedErr)
		})
	}
}
