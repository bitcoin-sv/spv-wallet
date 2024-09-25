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
)

func TestCreateTransactionDraftError(t *testing.T) {
	errorTests := map[string]struct {
		spec          *draft.TransactionSpec
		expectedError models.SPVError
	}{
		"return error for nil as transaction spec": {
			spec:          nil,
			expectedError: txerrors.ErrDraftSpecificationRequired,
		},
		"return error for transaction spec without xPub Id": {
			spec:          &draft.TransactionSpec{},
			expectedError: txerrors.ErrDraftSpecificationXPubIDRequired,
		},
		"return error for no outputs in transaction spec": {
			spec:          &draft.TransactionSpec{XPubID: fixtures.Sender.XPubID},
			expectedError: txerrors.ErrDraftRequiresAtLeastOneOutput,
		},
		"return error for empty output list in transaction spec": {
			spec: &draft.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(),
			},
			expectedError: txerrors.ErrDraftRequiresAtLeastOneOutput,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			draftService := given.NewDraftTransactionService()

			// when:
			tx, err := draftService.Create(context.Background(), test.spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}
