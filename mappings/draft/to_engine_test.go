package mappingsdraft_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	mappingsdraft "github.com/bitcoin-sv/spv-wallet/mappings/draft"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/stretchr/testify/require"
)

func TestMapToEngine(t *testing.T) {
	tests := map[string]struct {
		req      *request.DraftTransaction
		expected *draft.TransactionSpec
	}{
		"map op_return string output": {
			req: &request.DraftTransaction{
				Outputs: []request.Output{
					opreturn.Output{
						DataType: opreturn.DataTypeStrings,
						Data:     []string{"Example data"},
					},
				},
			},
			expected: &draft.TransactionSpec{
				Outputs: outputs.NewSpecifications(
					&outputs.OpReturn{
						DataType: opreturn.DataTypeStrings,
						Data:     []string{"Example data"},
					},
				),
			},
		},
		"map paymail output": {
			req: &request.DraftTransaction{
				Outputs: []request.Output{
					paymailreq.Output{
						To:       "receiver@example.com",
						Satoshis: 1000,
						From:     optional.Of("sender@example.com"),
					},
				},
			},
			expected: &draft.TransactionSpec{
				Outputs: outputs.NewSpecifications(
					&outputs.Paymail{
						To:       "receiver@example.com",
						Satoshis: 1000,
						From:     optional.Of("sender@example.com"),
					},
				),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			result, err := mappingsdraft.ToEngine(test.req)

			// then:
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		})
	}
}
