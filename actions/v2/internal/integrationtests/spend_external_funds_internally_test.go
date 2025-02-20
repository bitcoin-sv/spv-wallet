package integrationtests

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
)

func TestSpendExternalFundsInternally(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	alice := fixtures.Sender
	bob := fixtures.RecipientInternal

	// and:
	given.Paymail().ExternalPaymailHost().WillRespondWithP2PWithBEEFCapabilities()

	// and:
	externalTxReference := "z1cde5fa-7b29-403e-8a16-d92f7304b8c2"

	// when:
	receiveTx := given.TransactionScenario(alice).ReceivesFromExternal(10, externalTxReference)

	// then:
	then.User(alice).Balance().IsEqualTo(10)
	then.User(alice).Operations().Last().
		WithTxID(receiveTx.TxID).
		WithTxStatus("BROADCASTED").
		WithValue(10).
		WithType("incoming")

	// when:
	internalTx := given.TransactionScenario(alice).SendsToInternal(bob, 5)

	// then:
	then.User(alice).Balance().IsEqualTo(5)
	then.User(bob).Balance().IsEqualTo(5)

	then.User(alice).Operations().Last().
		WithTxID(internalTx.TxID).
		WithTxStatus("BROADCASTED").
		WithValue(-5).
		WithType("outgoing").
		WithCounterparty(bob.DefaultPaymail().Address())
}
