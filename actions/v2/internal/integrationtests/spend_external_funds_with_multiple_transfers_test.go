package integrationtests

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testsuite"
)

func TestSpendInternalFundsWithMultipleTransfers(t *testing.T) {
	testsuite.RunOnAllDBMS(t, func(t *testing.T, dbms string) {
		// given:
		given, when, then, cleanup := testsuite.SetupDBMSTest(t, dbms)
		defer cleanup()

		// and:
		when.Alice().ReceivesFromExternal(50)

		// then:
		then.Alice().Balance().IsEqualTo(50)

		// when:
		firstTxID := when.Alice().SendsFundsTo(given.Bob(), 30)

		// then:
		then.Alice().Balance().IsEqualTo(19)
		then.Bob().Balance().IsEqualTo(30)

		// and:
		then.Alice().Operations().Last().
			WithTxID(firstTxID).
			WithTxStatus("BROADCASTED").
			WithValue(-31).
			WithType("outgoing").
			WithCounterparty(given.Bob().DefaultPaymail().Address())

		// and:
		then.Bob().Operations().Last().
			WithTxID(firstTxID).
			WithTxStatus("BROADCASTED").
			WithValue(30).
			WithType("incoming").
			WithCounterparty(given.Alice().DefaultPaymail().Address())

		// when:
		secondTxID := when.Bob().SendsFundsTo(given.Charlie(), 20)

		// then:
		then.Bob().Balance().IsEqualTo(9)

		// and:
		then.Bob().Operations().Last().
			WithTxID(secondTxID).
			WithTxStatus("BROADCASTED").
			WithValue(-21).
			WithType("outgoing").
			WithCounterparty(given.Charlie().DefaultPaymail().Address())
	})
}
