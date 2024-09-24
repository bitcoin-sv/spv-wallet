package draft_test

import (
	"context"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePaymailDraft(t *testing.T) {
	const transactionSatoshiValue = bsv.Satoshis(1)
	var recipient = fixtures.RecipientExternal.DefaultPaymail()
	var sender = fixtures.Sender.DefaultPaymail()

	t.Run("return draft with payment to valid paymail address", func(t *testing.T) {
		given := testabilities.Given(t)

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
		require.NoError(t, err)
		require.NotNil(t, draftTx)

		// and:
		annotations := draftTx.Annotations
		require.Len(t, annotations.Outputs, 1)

		// and:
		annotation := annotations.Outputs[0]
		require.Equal(t, transaction.BucketBSV, annotation.Bucket)
		require.NotNil(t, annotation.Paymail)
		assert.Equal(t, recipient, annotation.Paymail.Receiver)
		assert.Equal(t, paymailHostResponse.Reference, annotation.Paymail.Reference)
		assert.Equal(t, sender, annotation.Paymail.Sender)

		// debug:
		t.Logf("BEEF: %s", draftTx.BEEF)

		// when:
		tx, err := sdk.NewTransactionFromBEEFHex(draftTx.BEEF)

		// then:
		require.NoError(t, err)
		require.Len(t, tx.Outputs, 1)

		// and:
		output := tx.Outputs[0]
		require.EqualValues(t, transactionSatoshiValue, output.Satoshis)
		require.Equal(t, paymailHostResponse.Scripts[0], output.LockingScriptHex())
	})

	t.Run("return draft with payment with multiple outputs from valid paymail address", func(t *testing.T) {
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
			HasOutputs(2).
			Output(0).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(firstOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[0]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference).
			And().
			Output(1).
			HasBucket(transaction.BucketBSV).
			HasSatoshis(secondOutputSatoshiValue).
			HasLockingScript(paymailHostResponse.Scripts[1]).
			IsPaymail().
			HasReceiver(recipient).
			HasSender(sender).
			HasReference(paymailHostResponse.Reference)
	})

	t.Run("return draft with default paymail in sender annotation", func(t *testing.T) {
		given := testabilities.Given(t)

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
		require.NoError(t, err)
		require.NotNil(t, draftTx)

		// and:
		annotation := draftTx.Annotations.Outputs[0]

		// then:
		assert.Equal(t, fixtures.UserWithMorePaymails.DefaultPaymail(), annotation.Paymail.Sender)
	})

	errorTests := map[string]struct {
		user          fixtures.User
		spec          *outputs.Paymail
		expectedError models.SPVError
	}{
		"for no paymail address": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"for only alias without domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "test",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"for domain without alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"for paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "$$$@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"for paymail with invalid domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       "test@example.com.$$$",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrReceiverPaymailAddressIsInvalid,
		},
		"for zero satoshis value": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: 0,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"for no satoshis value": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To: recipient,
			},
			expectedError: txerrors.ErrOutputValueTooLow,
		},
		"for sender paymail without domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for sender paymail without alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for sender paymail with invalid alias": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("$$$@example.com"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for sender paymail with invalid domain domain": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of("sender@example.com.$$$"),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for sender paymail address not existing in our system": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientExternal.DefaultPaymail()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for sender paymail not belonging to that user": {
			user: fixtures.Sender,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
				From:     optional.Of(fixtures.RecipientInternal.DefaultPaymail()),
			},
			expectedError: txerrors.ErrSenderPaymailAddressIsInvalid,
		},
		"for default sender paymail of user without paymail": {
			user: fixtures.UserWithoutPaymail,
			spec: &outputs.Paymail{
				To:       recipient,
				Satoshis: transactionSatoshiValue,
			},
			expectedError: txerrors.ErrDraftSenderPaymailAddressNoDefault,
		},
	}
	for name, test := range errorTests {
		t.Run("return error "+name, func(t *testing.T) {
			given := testabilities.Given(t)

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
			require.Error(t, err)
			require.ErrorIs(t, err, test.expectedError)
			require.Nil(t, tx)
		})
	}

	paymailErrorTests := map[string]struct {
		paymailHostScenario func(tpaymail.PaymailHostGiven)
		expectedError       models.SPVError
	}{
		"paymail host is responding with not found on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithNotFoundOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"paymail host is failing on capabilities": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithErrorOnCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"paymail host is not supporting p2p destinations capability": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithBasicCapabilities()
			},
			expectedError: pmerrors.ErrPaymailHostNotSupportingP2P,
		},
		"paymail host is failing on p2p destinations": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithErrorOnP2PDestinations()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"paymail host p2p destinations is returning not found": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithNotFoundOnP2PDestination()
			},
			expectedError: pmerrors.ErrPaymailHostResponseError,
		},
		"paymail host p2p destinations is responding with requirement for more sats then requested": {
			paymailHostScenario: func(paymailHost tpaymail.PaymailHostGiven) {
				paymailHost.WillRespondWithP2PDestinationsWithSats(transactionSatoshiValue + 1)
			},
			expectedError: pmerrors.ErrPaymailHostInvalidResponse,
		},
	}
	for name, test := range paymailErrorTests {
		t.Run("return error when "+name, func(t *testing.T) {
			given := testabilities.Given(t)

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
			require.Error(t, err)
			require.ErrorIs(t, err, test.expectedError)
			require.Nil(t, tx)
		})
	}
}
