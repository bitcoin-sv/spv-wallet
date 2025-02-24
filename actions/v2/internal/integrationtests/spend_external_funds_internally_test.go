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
	// TODO: started SPV Wallet V2 <- dodatkowa funkcja .StartedSPVWalletV2
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	// TODO: add users as aliases (vars) for alice, bob, charlie
	alice := fixtures.Sender
	bob := fixtures.RecipientInternal

	// and:
	given.Paymail().ExternalPaymailHost().WillRespondWithP2PWithBEEFCapabilities()

	// and:
	// fix it -> take a look at incoming_paymail_tx_test.go, we need to go through all steps

	// when:
	// wrap "Alice" in testabilities wrapped with methods
	receiveTx := given.TransactionScenario(alice).ReceivesFromExternal(10)

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
	then.User(alice).Balance().IsEqualTo(5) // Probably - fee (1), think about expected fee variable
	then.User(bob).Balance().IsEqualTo(5)

	then.User(alice).Operations().Last().
		WithTxID(internalTx.TxID).
		WithTxStatus("BROADCASTED").
		WithValue(-5).
		WithType("outgoing").
		WithCounterparty(bob.DefaultPaymail().Address())

	// TODO: Check op for BOB
}
