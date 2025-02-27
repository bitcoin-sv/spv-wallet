package integrationtests

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	"testing"
)

func TestSpendInternalFundsWithMultipleTransfers(t *testing.T) {
	// given:
	given, when, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletV2()
	defer cleanup()

	// and:
	when.Alice().ReceivesFromExternal(50)

	// then:
	then.Alice().Balance().IsEqualTo(50)

	// when:
	firstTxID := when.Alice().SendsTo(given.Bob(), 30)

	// then:
	then.Alice().Balance().IsEqualTo(0)
	then.Bob().Balance().IsEqualTo(30)

	// and:
	then.Alice().Operations().Last().
		WithTxID(firstTxID).
		WithTxStatus("BROADCASTED").
		WithValue(-50).
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
	secondTxID := when.Bob().SendsTo(given.Charlie(), 20)

	// add opreturn

	// then:
	then.Bob().Balance().IsEqualTo(0)

	// and:
	then.Bob().Operations().Last().
		WithTxID(secondTxID).
		WithTxStatus("BROADCASTED").
		WithValue(-30).
		WithType("outgoing").
		WithCounterparty(given.Charlie().DefaultPaymail().Address())

}
