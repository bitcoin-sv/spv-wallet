package request_test

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/stretchr/testify/require"
)

func TestDraft_TransactionJSON(t *testing.T) {
	tests := map[string]struct {
		json  string
		draft *request.DraftTransaction
	}{
		"OP_RETURN output with single string": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": [ "hello world" ]
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeStrings,
						Data:     []string{"hello world"},
					},
				},
			},
		},
		"OP_RETURN output with multiple strings": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": [ "hello", "world" ]
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeStrings,
						Data:     []string{"hello", "world"},
					},
				},
			},
		},
		"OP_RETURN output with default data type": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "data": [ "hello world" ]
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeDefault,
						Data:     []string{"hello world"},
					},
				},
			},
		},
		"OP_RETURN output with hex": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": [ "68656c6c6f20776f726c64" ]
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeHexes,
						Data:     []string{hex.EncodeToString([]byte("hello world"))},
					},
				},
			},
		},
		"OP_RETURN output with multiple hex": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "hexes",
				  "data": [ "68656c6c6f", "20776f726c64" ]
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeHexes,
						Data:     []string{hex.EncodeToString([]byte("hello")), hex.EncodeToString([]byte(" world"))},
					},
				},
			},
		},
		"Paymail output without sender": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver@example.com",
				  "satoshis": 1000
				}
			  ]
			}`,
			draft: &request.DraftTransaction{
				Outputs: []request.Output{
					paymailreq.Output{
						To:       "receiver@example.com",
						Satoshis: 1000,
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run("draft from JSON: "+name, func(t *testing.T) {
			var draft *request.DraftTransaction
			err := json.Unmarshal([]byte(test.json), &draft)
			require.NoError(t, err)
			require.Equal(t, test.draft, draft)
		})
		t.Run("draft to JSON: "+name, func(t *testing.T) {
			data, err := json.Marshal(test.draft)
			require.NoError(t, err)
			jsonValue := string(data)
			require.JSONEq(t, test.json, jsonValue)
		})
	}
}

func TestDraft_TransactionJSONParsingErrors(t *testing.T) {
	tests := map[string]struct {
		json        string
		expectedErr string
	}{
		"Unsupported output type": {
			json: `{
			  "outputs": [
				{
				  "type": "unsupported"
				}
			  ]
			}`,
			expectedErr: "unsupported output type",
		},
		"OP_RETURN output with unknown data type": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "unknown",
				  "data": [ "hello world" ]
				}
			  ]
			}`,
			expectedErr: "invalid data type",
		},
		"OP_RETURN output with string instead of array as data": {
			json: `{
			  "outputs": [
				{
				  "type": "op_return",
				  "dataType": "strings",
				  "data": "hello world"
				}
			  ]
			}`,
			expectedErr: "json: cannot unmarshal",
		},
		"Paymail output with negative satoshis": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver@example.com",
				  "satoshis": -1
				}
			  ]
			}`,
			expectedErr: "json: cannot unmarshal",
		},
		"Paymail output with too high satoshis value": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver@example.com",
				  "satoshis": ` + getTooLargeSatsValueToParse() + `
				}
			  ]
			}`,
			expectedErr: "json: cannot unmarshal",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var draft *request.DraftTransaction
			err := json.Unmarshal([]byte(test.json), &draft)
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

func getTooLargeSatsValueToParse() string {
	maxSats := strconv.FormatUint(math.MaxUint64, 10)
	return maxSats + "0"
}

func TestDraft_TransactionJSONEncodingErrors(t *testing.T) {
	tests := map[string]struct {
		draft       *request.DraftTransaction
		expectedErr string
	}{
		"Unsupported output type": {
			draft: &request.DraftTransaction{
				Outputs: []request.Output{&unsupportedOutput{}},
			},
			expectedErr: "unsupported output type",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := json.Marshal(test.draft)
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

type unsupportedOutput struct{}

func (o *unsupportedOutput) GetType() string {
	return "unsupported"
}
