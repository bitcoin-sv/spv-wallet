package draft_test

import (
	"context"
	"testing"

	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	tpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/testabilities"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
)

func TestCreatePaymailDraft(t *testing.T) {
	const transactionSatoshiValue = bsv.Satoshis(1)
	var recipient = fixtures.RecipientExternal.DefaultPaymail()
	var sender = fixtures.Sender.DefaultPaymail()

	t.Run("return draft with payment to valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		paymailHostResponse := given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue)

		// and:
		draftService := given.NewDraftTransactionService()

		// and:
		spec := &draft.TransactionSpec{
			XPubID: fixtures.Sender.XPubID,
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		then.Created(draftTx).WithNoError(err).WithParseableBEEFHex().
			HasOutputs(1).
			Output(0).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(transactionSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[0]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("style 1 __ return draft with payment with multiple outputs from valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// and:
		paymailHostResponse := given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		// and:
		draftService := given.NewDraftTransactionService()

		// and:
		spec := &draft.TransactionSpec{
			XPubID: fixtures.Sender.XPubID,
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		thenTx := then.Created(draftTx).WithNoError(err).WithParseableBEEFHex()

		thenTx.HasOutputs(2)

		thenTx.Output(0).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(firstOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[0]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

		thenTx.Output(1).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(secondOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[1]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("style 2 __ return draft with payment with multiple outputs from valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// and:
		paymailHostResponse := given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		// and:
		draftService := given.NewDraftTransactionService()

		// and:
		spec := &draft.TransactionSpec{
			XPubID: fixtures.Sender.XPubID,
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		then.Created(draftTx).WithNoError(err).
			WithParseableBEEFHex().
			HasOutputs(2).

			// and:
			Output(0).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(firstOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[0]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference).
			And().

			// and:
			Output(1).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(secondOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[1]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)

	})

	t.Run("style 3 __ return draft with payment with multiple outputs from valid paymail address", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// and:
		paymailHostResponse := given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		// and:
		draftService := given.NewDraftTransactionService()

		// and:
		spec := &draft.TransactionSpec{
			XPubID: fixtures.Sender.XPubID,
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
				From:     optional.Of(sender),
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		then.Created(draftTx).WithNoError(err).
			WithParseableBEEFHex().
			HasOutputs(2).

			// and:
			HasOutput(0, func(that testabilities.OutputAssertion) {
				that.HasBucket(transaction.BucketBSV).
					HasSatoshis(firstOutputSatoshiValue).
					HasLockingScript(paymailHostResponse.Scripts[0]).
					IsPaymail().
					HasReceiver(recipient).
					HasSender(sender).
					HasReference(paymailHostResponse.Reference)
			}).
			HasOutput(1, func(that testabilities.OutputAssertion) {
				that.HasBucket(transaction.BucketBSV).
					HasSatoshis(secondOutputSatoshiValue).
					HasLockingScript(paymailHostResponse.Scripts[1]).
					IsPaymail().
					HasReceiver(recipient).
					HasSender(sender).
					HasReference(paymailHostResponse.Reference)
			})
	})

	t.Run("return draft with default paymail in sender annotation", func(t *testing.T) {
		given, then := testabilities.New(t)

		// given:
		given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue)

		// and:
		draftService := given.NewDraftTransactionService()

		// and:
		spec := &draft.TransactionSpec{
			XPubID: fixtures.UserWithMorePaymails.XPubID,
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		then.Created(draftTx).WithNoError(err).WithParseableBEEFHex().
			Output(0).
			IsPaymail().
			HasSender(fixtures.UserWithMorePaymails.DefaultPaymail())
	})

	errorTests := map[string]struct {
		user          fixtures.User
		spec          *outputs.Paymail
		expectedError models.SPVError
	}{
		"return error for no paymail address": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for only alias without domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "test",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for domain without alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "$$$@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for paymail with invalid domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "test@example.com.$$$",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"return error for zero satoshis value": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: 0,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for no satoshis value": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To: recipient,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"return error for sender paymail without domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail without alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("$$$@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail with invalid domain domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender@example.com.$$$"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail address not existing in our system": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientExternal.DefaultPaymail()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for sender paymail not belonging to that user": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientInternal.DefaultPaymail()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"return error for default sender paymail of user without paymail": {
			user: fixtures.UserWithoutPaymail,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrDraftSenderPaymailAddressNoDefault,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			given.ExternalRecipientHost().WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue)

			// and:
			draftService := given.NewDraftTransactionService()

			// and:
			spec := &draft.TransactionSpec{
				XPubID:  test.user.XPubID,
				Outputs: outputs.NewSpecifications(test.spec),
			}

			// when:
			tx, err := draftService.Create(context.Background(), spec)

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
			draftService := given.NewDraftTransactionService()

			// and:
			spec := &draft.TransactionSpec{
				XPubID: fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(&outputs.Paymail{
					To:       recipient,
					Satoshis: transactionSatoshiValue,
				}),
			}

			// when:
			tx, err := draftService.Create(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}
