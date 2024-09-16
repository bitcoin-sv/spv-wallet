package draft_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/stretchr/testify/require"
)

func TestCreateTransactionDraftError(t *testing.T) {
	errorTests := map[string]struct {
		spec          *draft.TransactionSpec
		expectedError string
	}{
		"for nil as transaction spec": {
			spec:          nil,
			expectedError: "draft requires a specification",
		},
		"for no outputs in transaction spec": {
			spec:          &draft.TransactionSpec{},
			expectedError: "draft requires at least one output",
		},
		"for empty output list in transaction spec": {
			spec: &draft.TransactionSpec{Outputs: outputs.NewSpecifications()},
		},
	}
	for name, test := range errorTests {
		t.Run("return error "+name, func(t *testing.T) {
			// given:
			draftService := draft.NewDraftService(paymailmock.CreatePaymailClientService("test"), tester.Logger())

			// when:
			_, err := draftService.Create(context.Background(), test.spec)

			// then:
			require.Error(t, err)
			require.ErrorContains(t, err, test.expectedError)
		})
	}
}
