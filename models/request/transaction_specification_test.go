package request_test

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/stretchr/testify/require"
)

func TestTransactionSpecification_TransactionJSON(t *testing.T) {
	tests := map[string]struct {
		json string
		spec *request.TransactionSpecification
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
			spec: &request.TransactionSpecification{
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
			spec: &request.TransactionSpecification{
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
			spec: &request.TransactionSpecification{
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
			spec: &request.TransactionSpecification{
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
			spec: &request.TransactionSpecification{
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
			spec: &request.TransactionSpecification{
				Outputs: []request.Output{
					paymailreq.Output{
						To:       "receiver@example.com",
						Satoshis: 1000,
					},
				},
			},
		},
		"Paymail output with sender": {
			json: `{
			  "outputs": [
				{
				  "type": "paymail",
				  "to": "receiver@example.com",
				  "satoshis": 1000,
				  "from": "sender@example.com"
				}
			  ]
			}`,
			spec: &request.TransactionSpecification{
				Outputs: []request.Output{
					paymailreq.Output{
						To:       "receiver@example.com",
						Satoshis: 1000,
						From:     optional.Of("sender@example.com"),
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run("spec from JSON: "+name, func(t *testing.T) {
			var spec *request.TransactionSpecification
			err := json.Unmarshal([]byte(test.json), &spec)
			require.NoError(t, err)
			require.Equal(t, test.spec, spec)
		})
		t.Run("spec to JSON: "+name, func(t *testing.T) {
			data, err := json.Marshal(test.spec)
			require.NoError(t, err)
			jsonValue := string(data)
			require.JSONEq(t, test.json, jsonValue)
		})
	}
}

func TestTransactionSpecification_JSONParsingErrors(t *testing.T) {
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
			var spec *request.TransactionSpecification
			err := json.Unmarshal([]byte(test.json), &spec)
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

func getTooLargeSatsValueToParse() string {
	maxSats := strconv.FormatUint(math.MaxUint64, 10)
	return maxSats + "0"
}

func TestTransactionSpecification_JSONEncodingErrors(t *testing.T) {
	tests := map[string]struct {
		spec        *request.TransactionSpecification
		expectedErr string
	}{
		"Unsupported output type": {
			spec: &request.TransactionSpecification{
				Outputs: []request.Output{&unsupportedOutput{}},
			},
			expectedErr: "unsupported output type",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := json.Marshal(test.spec)
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

type unsupportedOutput struct{}

func (o *unsupportedOutput) GetType() string {
	return "unsupported"
}
