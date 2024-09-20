package draft_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePaymailDraft(t *testing.T) {
	const testDomain = "example.com"
	const transactionSatoshiValue = bsv.Satoshis(1)
	const recipient = "test" + "@" + testDomain

	t.Run("return draft with payment to valid paymail address", func(t *testing.T) {
		// given:
		paymailHostResponse := paymailmock.P2PDestinationsForSats(transactionSatoshiValue)

		paymailClient := paymailmock.CreatePaymailClientService(testDomain)
		paymailClient.WillRespondWithP2PCapabilities()
		paymailClient.
			WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
			With(paymailHostResponse)

		// and:
		draftService := draft.NewDraftService(paymailClient, tester.Logger())

		// and:
		spec := &draft.TransactionSpec{
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
		annotations := draftTx.Annotations
		require.Len(t, annotations.Outputs, 1)

		// and:
		annotation := annotations.Outputs[0]
		require.Equal(t, transaction.BucketBSV, annotation.Bucket)
		require.NotNil(t, annotation.Paymail)
		assert.Equal(t, recipient, annotation.Paymail.Receiver)
		assert.Equal(t, paymailHostResponse.Reference, annotation.Paymail.Reference)

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
		const firstOutputSatoshiValue = bsv.Satoshis(1)
		const secondOutputSatoshiValue = bsv.Satoshis(2)
		const paymentSatoshiValue = firstOutputSatoshiValue + secondOutputSatoshiValue

		// given:
		paymailHostResponse := paymailmock.P2PDestinationsForSats(firstOutputSatoshiValue, secondOutputSatoshiValue)

		paymailClient := paymailmock.CreatePaymailClientService(testDomain)
		paymailClient.WillRespondWithP2PCapabilities()
		paymailClient.
			WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
			With(paymailHostResponse)

		// and:
		draftService := draft.NewDraftService(paymailClient, tester.Logger())

		// and:
		spec := &draft.TransactionSpec{
			Outputs: outputs.NewSpecifications(&outputs.Paymail{
				To:       recipient,
				Satoshis: paymentSatoshiValue,
			}),
		}

		// when:
		draftTx, err := draftService.Create(context.Background(), spec)

		// then:
		require.NoError(t, err)
		require.NotNil(t, draftTx)

		// and:
		annotations := draftTx.Annotations
		require.Len(t, annotations.Outputs, 2)
		assert.Equal(t, transaction.BucketBSV, annotations.Outputs[0].Bucket)
		assert.Equal(t, transaction.BucketBSV, annotations.Outputs[1].Bucket)
		// TODO: add assertions for paymail annotations

		// debug:
		t.Logf("BEEF: %s", draftTx.BEEF)

		// when:
		tx, err := sdk.NewTransactionFromBEEFHex(draftTx.BEEF)

		// then:
		require.NoError(t, err)
		require.Len(t, tx.Outputs, 2)
		require.EqualValues(t, firstOutputSatoshiValue, tx.Outputs[0].Satoshis)
		require.Equal(t, paymailHostResponse.Scripts[0], tx.Outputs[0].LockingScriptHex())
		require.EqualValues(t, secondOutputSatoshiValue, tx.Outputs[1].Satoshis)
		require.Equal(t, paymailHostResponse.Scripts[1], tx.Outputs[1].LockingScriptHex())
	})

	errorTests := map[string]struct {
		spec          *outputs.Paymail
		expectedError string
	}{
		"for only alias without domain": {
			spec: &outputs.Paymail{
				To:       "test",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: "paymail address is invalid",
		},
		"for domain without alias": {
			spec: &outputs.Paymail{
				To:       "@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: "paymail address is invalid",
		},
		"for paymail with invalid alias": {
			spec: &outputs.Paymail{
				To:       "$$$@example.com",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: "paymail address is invalid",
		},
		"for paymail with invalid domain": {
			spec: &outputs.Paymail{
				To:       "test@example.com.$$$",
				Satoshis: transactionSatoshiValue,
			},
			expectedError: "paymail address is invalid",
		},
	}
	for name, test := range errorTests {
		t.Run("return error "+name, func(t *testing.T) {
			// given:
			paymailClient := paymailmock.CreatePaymailClientService(testDomain)
			paymailClient.WillRespondWithP2PCapabilities()

			// and:
			draftService := draft.NewDraftService(paymailClient, tester.Logger())

			// and:
			spec := &draft.TransactionSpec{
				Outputs: outputs.NewSpecifications(test.spec),
			}

			// when:
			tx, err := draftService.Create(context.Background(), spec)

			// then:
			require.Error(t, err)
			require.ErrorContains(t, err, test.expectedError)
			require.Nil(t, tx)
		})
	}

	paymailErrorTests := map[string]struct {
		paymailHostScenario func(*paymailmock.PaymailClientMock)
		expectedError       string
	}{
		"paymail host is responding with not found on capabilities": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.WillRespondWithNotFoundOnCapabilities()
			},
			expectedError: "paymail host is responding with error",
		},
		"paymail host is failing on capabilities": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.WillRespondWithErrorOnCapabilities()
			},
			expectedError: "paymail host is responding with error",
		},
		"paymail host is not supporting p2p destinations capability": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.WillRespondWithBasicCapabilities()
			},
			expectedError: "paymail host is not supporting P2P capabilities",
		},
		"paymail host is failing on p2p destinations": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.
					WillRespondWithP2PCapabilities().
					WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
					WithInternalServerError()
			},
			expectedError: "paymail host is responding with error",
		},
		"paymail host p2p destinations is returning not found": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.
					WillRespondWithP2PCapabilities().
					WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
					WithNotFound()
			},
		},
		"paymail host p2p destinations is responding with requirement for more sats then requested": {
			paymailHostScenario: func(paymailClient *paymailmock.PaymailClientMock) {
				paymailClient.
					WillRespondWithP2PCapabilities().
					WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).
					With(paymailmock.P2PDestinationsForSats(transactionSatoshiValue + 1))
			},
			expectedError: "paymail host invalid response",
		},
	}
	for name, test := range paymailErrorTests {
		t.Run("return error when "+name, func(t *testing.T) {
			// given:
			paymailClient := paymailmock.CreatePaymailClientService(testDomain)
			test.paymailHostScenario(paymailClient.PaymailClientMock)

			// and:
			draftService := draft.NewDraftService(paymailClient, tester.Logger())

			// and:
			spec := &draft.TransactionSpec{
				Outputs: outputs.NewSpecifications(&outputs.Paymail{
					To:       recipient,
					Satoshis: transactionSatoshiValue,
				}),
			}

			// when:
			tx, err := draftService.Create(context.Background(), spec)

			// then:
			require.Error(t, err)
			require.ErrorContains(t, err, test.expectedError)
			require.Nil(t, tx)
		})
	}
}
