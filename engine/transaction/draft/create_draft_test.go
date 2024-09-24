package draft_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/testabilities"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestCreateTransactionDraftError(t *testing.T) {
	errorTests := map[string]struct {
		spec          *draft.TransactionSpec
		expectedError models.SPVError
	}{
		"for nil as transaction spec": {
			spec:          nil,
			expectedError: txerrors.ErrDraftSpecificationRequired,
		},
		"for transaction spec without xPub Id": {
			spec:          &draft.TransactionSpec{},
			expectedError: txerrors.ErrDraftSpecificationXPubIDRequired,
		},
		"for no outputs in transaction spec": {
			spec:          &draft.TransactionSpec{XPubID: fixtures.Sender.XPubID},
			expectedError: txerrors.ErrDraftRequiresAtLeastOneOutput,
		},
		"for empty output list in transaction spec": {
			spec: &draft.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(),
			},
			expectedError: txerrors.ErrDraftRequiresAtLeastOneOutput,
		},
	}
	for name, test := range errorTests {
		t.Run("return error "+name, func(t *testing.T) {
			given := testabilities.Given(t)

			// given:
			draftService := given.NewDraftTransactionService()

			// when:
			_, err := draftService.Create(context.Background(), test.spec)

			// then:
			require.Error(t, err)
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}
