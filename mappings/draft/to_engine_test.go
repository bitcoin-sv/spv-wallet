package mappingsdraft_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	mappingsdraft "github.com/bitcoin-sv/spv-wallet/mappings/draft"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
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
					&opreturn.Output{
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
