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

func TestCreateTransactionOutlineFromMinimumSpec(t *testing.T) {
	t.Run("minimum beef transaction spec", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		service := given.NewTransactionOutlinesService()

		// when:
		tx, err := service.CreateBEEF(context.Background(), given.MinimumValidTransactionSpec())

		// then:
		then.Created(tx).WithNoError(err).WithParseableBEEFHex()
	})

	t.Run("minimum raw transaction spec", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		service := given.NewTransactionOutlinesService()

		// when:
		tx, err := service.CreateRawTx(context.Background(), given.MinimumValidTransactionSpec())

		// then:
		then.Created(tx).WithNoError(err).WithParseableRawHex()
	})
}

func TestCreateTransactionOutlineBEEFError(t *testing.T) {
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
			tx, err := service.CreateBEEF(context.Background(), test.spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func TestCreateTransactionOutlineRAWError(t *testing.T) {
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
			tx, err := service.CreateRawTx(context.Background(), test.spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}

	t.Run("return error when user has not enough funds", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		service := given.NewTransactionOutlinesService()

		// and:
		given.UserHasNotEnoughFunds()

		// when:
		tx, err := service.CreateRawTx(context.Background(), given.MinimumValidTransactionSpec())

		// then:
		then.Created(tx).WithError(err).ThatIs(txerrors.ErrTxOutlineInsufficientFunds)

	})
}
