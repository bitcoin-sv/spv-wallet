package integrationtests

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	"testing"
)

func TestSpendExternalFundsOpReturn(t *testing.T) {
	// given:
	given, when, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletV2()
	defer cleanup()

	// and:
	receiveTxID := when.Alice().ReceivesFromExternal(2)

	// then:
	then.Alice().Balance().IsEqualTo(2)
	then.Alice().Operations().Last().
		WithTxID(receiveTxID).
		WithTxStatus("BROADCASTED").
		WithValue(2).
		WithType("incoming")

	// when:
	internalTxID := when.Alice().SendsData([]string{"Hello", "Bob!"})

	// then:
	then.Alice().Balance().IsEqualTo(0)

	then.Alice().Operations().Last().
		WithTxID(internalTxID).
		WithTxStatus("BROADCASTED").
		WithValue(-2).
		WithType("outgoing")

}
