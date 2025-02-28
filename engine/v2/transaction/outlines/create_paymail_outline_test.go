package outlines_test

import (
	"context"
	"testing"

	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func TestCreatePaymailTransactionOutlineBEEF(t *testing.T) {
	const transactionSatoshiValue = bsv.Satoshis(1)
	var recipient = fixtures.RecipientExternal.DefaultPaymail().Address()
	var sender = fixtures.Sender.DefaultPaymail().Address()

	t.Run("return transaction outline with payment to valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		tx, err := service.CreateBEEF(context.Background(), spec)

		// then:
		paymailHostResponse := then.ExternalPaymailHost().ReceivedP2PDestinationRequest(transactionSatoshiValue)

		thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

		thenTx.HasOutputs(1)

		thenTx.Output(0).
			HasBucket(bucket.BSV).
			HasSatoshis(transactionSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("return transaction outline with payment with multiple outputs from valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)
		satoshisToSplit := bsv.Satoshis(33)
		var splits uint64 = 3
		expectedSplitValue := bsv.Satoshis(11)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: satoshisToSplit,
				Splits:   splits,
				From:     optional.Of(sender),
			}),
		}

		// when:
		tx, err := service.CreateBEEF(context.Background(), spec)

		// then:
		paymailHostResponse := then.ExternalPaymailHost().ReceivedP2PDestinationRequest(satoshisToSplit)

		thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

		thenTx.HasOutputs(3)

		thenTx.Output(0).
			HasBucket(bucket.BSV).
			HasSatoshis(expectedSplitValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

		thenTx.Output(1).
			HasBucket(bucket.BSV).
			HasSatoshis(expectedSplitValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

		thenTx.Output(2).
			HasBucket(bucket.BSV).
			HasSatoshis(expectedSplitValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

	})

	t.Run("return transaction outline with payment split in multiple outputs", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// and:
		given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		tx, err := service.CreateBEEF(context.Background(), spec)

		// then:
		paymailHostResponse := then.ExternalPaymailHost().ReceivedP2PDestinationRequest(paymentSatoshiValue)

		thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

		thenTx.HasOutputs(2)

		thenTx.Output(0).
			HasBucket(bucket.BSV).
			HasSatoshis(firstOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

		thenTx.Output(1).
			HasBucket(bucket.BSV).
			HasSatoshis(secondOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[1].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("return transaction outline with default paymail in sender annotation", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.UserWithMorePaymails.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			}),
		}

		// when:
		tx, err := service.CreateBEEF(context.Background(), spec)

		// then:
		then.ExternalPaymailHost().ReceivedP2PDestinationRequest(transactionSatoshiValue)

		// and:
		then.Created(tx).WithNoError(err).WithParseableBEEFHex().
			Output(0).
			IsPaymail().
			HasSender(fixtures.UserWithMorePaymails.DefaultPaymail().Address())
	})

	errorTests := map[string]struct {
		user          fixtures.User
		spec          *outlines.Paymail
		expectedError models.SPVError
	}{
		"return error for no paymail address": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for only alias without domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "test",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for domain without alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "$$$@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "test@example.com.$$$",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for zero satoshis value": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: 0,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for no satoshis value": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To: recipient,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for sender paymail without domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail without alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("$$$@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid domain domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender@example.com.$$$"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail address not existing in our system": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientExternal.DefaultPaymail().Address()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail not belonging to that user": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientInternal.DefaultPaymail().Address()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for default sender paymail of user without paymail": {
			user: fixtures.UserWithoutPaymail,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrTxOutlineSenderPaymailAddressNoDefault,
		},
		"return error when satoshis are not divisible by splits": {
			user: fixtures.UserWithoutPaymail,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: 7,
				Splits:   3,
			},
			expectedError: txerrors.ErrTxOutlinePaymailSatoshisMustBeDivisibleBySplits,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

			// and:
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				UserID:  test.user.ID(),
				Outputs: outlines.NewOutputsSpecs(test.spec),
			}

			// when:
			tx, err := service.CreateBEEF(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}

	t.Run("return error when want to split paymail payment but recipient split it by himself", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(2, 4)

		// given:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: 6,
				Splits:   2,
			}),
		}

		// when:
		tx, err := service.CreateBEEF(context.Background(), spec)

		// then:
		then.Created(tx).WithError(err).ThatIs(txerrors.ErrTxOutlinePaymailCannotSplitWhenRecipientSplitting)
	})

	paymailErrorTests := map[string]struct {
		paymailHostScenario func(tpaymail.PaymailHostFixture)
		expectedError       models.SPVError
	}{
		"return error when paymail host is responding with not found on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithNotFoundOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host is failing on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithErrorOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host is not supporting p2p destinations capability": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithBasicCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostNotSupportingP2P,
		},
		"return error when paymail host is failing on p2p destinations": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithErrorOnP2PDestinations()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host p2p destinations is returning not found": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithNotFoundOnP2PDestination()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host p2p destinations is responding with requirement for more sats then requested": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue + 1)
			},
			expectedError: pmerrors.ErrPaymailHostInvalidResponse,
		},
	}
	for name, test := range paymailErrorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			test.paymailHostScenario(given.ExternalRecipientHost())

			// given:
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				UserID: fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
					To:       recipient,
					Satoshis: transactionSatoshiValue,
				}),
			}

			// when:
			tx, err := service.CreateBEEF(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func TestCreatePaymailTransactionOutlineRAW(t *testing.T) {
	const transactionSatoshiValue = bsv.Satoshis(1)
	var recipient = fixtures.RecipientExternal.DefaultPaymail().Address()
	var sender = fixtures.Sender.DefaultPaymail().Address()

	t.Run("return transaction outline with payment to valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		tx, err := service.CreateRawTx(context.Background(), spec)

		// then:
		paymailHostResponse := then.ExternalPaymailHost().ReceivedP2PDestinationRequest(transactionSatoshiValue)

		// and:
		thenTx := then.Created(tx).WithNoError(err).WithParseableRawHex()

		thenTx.IsWithoutTimeLock()

		thenTx.HasInputs(1)

		thenTx.Input(0).
			HasOutpoint(testabilities.UserFundsTransactionOutpoint).
			HasCustomInstructions(testabilities.UserFundsTransactionCustomInstructions)

		thenTx.HasOutputs(1)

		thenTx.Output(0).
			HasBucket(bucket.BSV).
			HasSatoshis(transactionSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("return transaction outline with payment with multiple outputs from valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// and:
		given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.Sender.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		tx, err := service.CreateRawTx(context.Background(), spec)

		// then:
		paymailHostResponse := then.ExternalPaymailHost().ReceivedP2PDestinationRequest(paymentSatoshiValue)

		// and:
		thenTx := then.Created(tx).WithNoError(err).WithParseableRawHex()

		thenTx.IsWithoutTimeLock()

		thenTx.HasInputs(1)

		thenTx.Input(0).
			HasOutpoint(testabilities.UserFundsTransactionOutpoint).
			HasCustomInstructions(testabilities.UserFundsTransactionCustomInstructions)

		thenTx.HasOutputs(2)

		thenTx.Output(0).
			HasBucket(bucket.BSV).
			HasSatoshis(firstOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[0].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

		thenTx.Output(1).
			HasBucket(bucket.BSV).
			HasSatoshis(secondOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Outputs[1].Script).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("return transaction outline with default paymail in sender annotation", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

		// and:
		service := given.NewTransactionOutlinesService()

		// and:
		spec := &outlines.TransactionSpec{
			UserID: fixtures.UserWithMorePaymails.ID(),
			Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			}),
		}

		// when:
		tx, err := service.CreateRawTx(context.Background(), spec)

		// then:
		thenTx := then.Created(tx).WithNoError(err).WithParseableRawHex()

		thenTx.IsWithoutTimeLock()

		thenTx.HasInputs(1)

		thenTx.Input(0).
			HasOutpoint(testabilities.UserFundsTransactionOutpoint).
			HasCustomInstructions(testabilities.UserFundsTransactionCustomInstructions)

		thenTx.Output(0).
			IsPaymail().
			HasSender(fixtures.UserWithMorePaymails.DefaultPaymail().Address())
	})

	errorTests := map[string]struct {
		user          fixtures.User
		spec          *outlines.Paymail
		expectedError models.SPVError
	}{
		"return error for no paymail address": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for only alias without domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "test",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for domain without alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "$$$@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       "test@example.com.$$$",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for zero satoshis value": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: 0,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for no satoshis value": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To: recipient,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for sender paymail without domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail without alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("$$$@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid domain domain": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender@example.com.$$$"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail address not existing in our system": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientExternal.DefaultPaymail().Address()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail not belonging to that user": {
			user: fixtures.Sender,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientInternal.DefaultPaymail().Address()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for default sender paymail of user without paymail": {
			user: fixtures.UserWithoutPaymail,
			spec: &outlines.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrTxOutlineSenderPaymailAddressNoDefault,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			given.ExternalRecipientHost().WillRespondWithP2PCapabilities()

			// and:
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				UserID:  test.user.ID(),
				Outputs: outlines.NewOutputsSpecs(test.spec),
			}

			// when:
			tx, err := service.CreateRawTx(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}

	paymailErrorTests := map[string]struct {
		paymailHostScenario func(tpaymail.PaymailHostFixture)
		expectedError       models.SPVError
	}{
		"return error when paymail host is responding with not found on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithNotFoundOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host is failing on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithErrorOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host is not supporting p2p destinations capability": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithBasicCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostNotSupportingP2P,
		},
		"return error when paymail host is failing on p2p destinations": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithErrorOnP2PDestinations()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host p2p destinations is returning not found": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithNotFoundOnP2PDestination()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"return error when paymail host p2p destinations is responding with requirement for more sats then requested": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostFixture) {
				paymailHost.WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue + 1)
			},
			expectedError: pmerrors.ErrPaymailHostInvalidResponse,
		},
	}
	for name, test := range paymailErrorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			test.paymailHostScenario(given.ExternalRecipientHost())

			// given:
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				UserID: fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(&outlines.Paymail{
					To:       recipient,
					Satoshis: transactionSatoshiValue,
				}),
			}

			// when:
			tx, err := service.CreateRawTx(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}
