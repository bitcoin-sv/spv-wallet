package request

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	"github.com/stretchr/testify/require"
)

func TestDraftTransactionJSON(t *testing.T) {
	tests := map[string]struct {
		json  string
		draft *DraftTransaction
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
			draft: &DraftTransaction{
				Outputs: []Output{
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
			draft: &DraftTransaction{
				Outputs: []Output{
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
			draft: &DraftTransaction{
				Outputs: []Output{
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
			draft: &DraftTransaction{
				Outputs: []Output{
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
			draft: &DraftTransaction{
				Outputs: []Output{
					opreturn.Output{
						DataType: opreturn.DataTypeHexes,
						Data:     []string{hex.EncodeToString([]byte("hello")), hex.EncodeToString([]byte(" world"))},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run("draft from JSON: "+name, func(t *testing.T) {
			var draft *DraftTransaction
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

func TestDraftTransactionJSONParsingErrors(t *testing.T) {
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
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var draft *DraftTransaction
			err := json.Unmarshal([]byte(test.json), &draft)
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

func TestDraftTransactionJSONEncodingErrors(t *testing.T) {
	tests := map[string]struct {
		draft       *DraftTransaction
		expectedErr string
	}{
		"Unsupported output type": {
			draft: &DraftTransaction{
				Outputs: []Output{&unsupportedOutput{}},
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
