package outlines_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func TestCreateTransactionOutlineError(t *testing.T) {
	errorTests := map[string]struct {
		spec          *outlines.TransactionSpec
		expectedError models.SPVError
	}{
		"return error for nil as transaction spec": {
			spec:          nil,
			expectedError: txerrors.ErrTxOutlineSpecificationRequired,
		},
		"return error for transaction spec without xPub Id": {
			spec:          &outlines.TransactionSpec{},
			expectedError: txerrors.ErrTxOutlineSpecificationUserIDRequired,
		},
		"return error for no outputs in transaction spec": {
			spec:          &outlines.TransactionSpec{UserID: fixtures.Sender.ID()},
			expectedError: txerrors.ErrTxOutlineRequiresAtLeastOneOutput,
		},
		"return error for empty output list in transaction spec": {
			spec: &outlines.TransactionSpec{
				UserID:  fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(),
			},
			expectedError: txerrors.ErrTxOutlineRequiresAtLeastOneOutput,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			service := given.NewTransactionOutlinesService()

			// when:
			tx, err := service.Create(context.Background(), test.spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}
