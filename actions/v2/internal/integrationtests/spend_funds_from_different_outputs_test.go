package integrationtests

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	"testing"
)

func TestSpendingFromMultipleSourceOutputs(t *testing.T) {
	// given
	given, when, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletV2()
	defer cleanup()

	// and:
	_ = when.Alice().ReceivesFromExternal(15)
	_ = when.Alice().ReceivesFromExternal(25)
	thirdTx := when.Alice().ReceivesFromExternal(10)

	// then:
	then.Alice().Balance().IsEqualTo(50) // 15 + 25 + 10

	// and:
	then.Alice().Operations().Last().
		WithTxID(thirdTx).
		WithTxStatus("BROADCASTED").
		WithValue(10).
		WithType("incoming")

	// when:
	txID4 := when.Alice().SendsFundsTo(given.Bob(), 41)

	// then:
	then.Alice().Balance().IsEqualTo(0)
	then.Bob().Balance().IsEqualTo(41)

	// and:
	then.Alice().Operations().Last().
		WithTxID(txID4).
		WithTxStatus("BROADCASTED").
		WithValue(-50).
		WithType("outgoing").
		WithCounterparty(given.Bob().DefaultPaymail().Address())

	// and:
	then.Bob().Operations().Last().
		WithTxID(txID4).
		WithTxStatus("BROADCASTED").
		WithValue(41).
		WithType("incoming").
		WithCounterparty(given.Alice().DefaultPaymail().Address())
}
