package integrationtests

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
)

func TestSpendExternalFundsInternally(t *testing.T) {
	// given:
	given, when, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletV2()
	defer cleanup()

	// and:
	receiveTxID := when.Alice().ReceivesFromExternal(10)

	// then:
	then.Alice().Balance().IsEqualTo(10)
	then.Alice().Operations().Last().
		WithTxID(receiveTxID).
		WithTxStatus("BROADCASTED").
		WithValue(10).
		WithType("incoming")

	// when:
	internalTxID := when.Alice().SendsTo(given.Bob(), 5)

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
}
