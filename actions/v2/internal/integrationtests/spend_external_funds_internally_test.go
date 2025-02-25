package integrationtests

import (
	internaltestabilities "github.com/bitcoin-sv/spv-wallet/actions/v2/internal/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"testing"
)

func TestSpendExternalFundsInternally(t *testing.T) {
	// given:
	tc := internaltestabilities.NewActorTests(t)
	cleanup := tc.Given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	tc.Given.Paymail().ExternalPaymailHost().WillRespondWithP2PWithBEEFCapabilities()

	// when:
	receiveTxID := tc.Alice.ReceivesFromExternal(10)

	// then:
	tc.Then(tc.Alice).Balance().IsEqualTo(10)
	tc.Then(tc.Alice).Operations().Last().
		WithTxID(receiveTxID).
		WithTxStatus("BROADCASTED").
		WithValue(10).
		WithType("incoming")

	// when:
	internalTxID := tc.Alice.SendsTo(tc.Bob, 5)

	// then:
	tc.Then(tc.Alice).Balance().IsEqualTo(0)
	tc.Then(tc.Bob).Balance().IsEqualTo(5)

	tc.Then(tc.Alice).Operations().Last().
		WithTxID(internalTxID).
		WithTxStatus("BROADCASTED").
		WithValue(-10).
		WithType("outgoing").
		WithCounterparty(tc.Bob.DefaultPaymail().Address())

	tc.Then(tc.Bob).Operations().Last().
		WithTxID(internalTxID).
		WithTxStatus("BROADCASTED").
		WithValue(5).
		WithType("incoming").
		WithCounterparty(tc.Alice.DefaultPaymail().Address())
}
