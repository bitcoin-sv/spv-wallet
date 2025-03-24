package integrationtests

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testsuite"
)

func TestSpendExternalFundsInternally(t *testing.T) {
	testsuite.RunOnAllDBMS(t, func(t *testing.T, dbms string) {
		// given:
		given, when, then, cleanup := testsuite.SetupDBMSTest(t, dbms)
		defer cleanup()

		// when:
		receiveTxID := when.Alice().ReceivesFromExternal(10)

		// then:
		then.Alice().Balance().IsEqualTo(10)
		then.Alice().Operations().Last().
			WithTxID(receiveTxID).
			WithTxStatus("BROADCASTED").
			WithValue(10).
			WithType("incoming")

		// when:
		internalTxID := when.Alice().SendsFundsTo(given.Bob(), 5)

		// then:
		then.Alice().Balance().IsEqualTo(4)
		then.Bob().Balance().IsEqualTo(5)

		then.Alice().Operations().Last().
			WithTxID(internalTxID).
			WithTxStatus("BROADCASTED").
			WithValue(-6).
			WithType("outgoing").
			WithCounterparty(given.Bob().DefaultPaymail().Address())

		then.Bob().Operations().Last().
			WithTxID(internalTxID).
			WithTxStatus("BROADCASTED").
			WithValue(5).
			WithType("incoming").
			WithCounterparty(given.Alice().DefaultPaymail().Address())
	})
}
