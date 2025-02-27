package integrationtests

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	"testing"
)

func TestSpendExternalFundsWithMultipleOutputs(t *testing.T) {
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
	// move to two methods for example OpReturn Paymail etc
	dataOutput := map[string]any{
		"type":     "op_return",
		"dataType": "strings",
		"data":     []string{"Hello, world!"},
	}

	internalTxID := when.Alice().SendsTo(given.Bob(), 1, dataOutput)

	// then:
	then.Alice().Balance().IsEqualTo(0)
	then.Bob().Balance().IsEqualTo(1)

	then.Alice().Operations().Last().
		WithTxID(internalTxID).
		WithTxStatus("BROADCASTED").
		WithValue(-2).
		WithType("outgoing").
		WithCounterparty(given.Bob().DefaultPaymail().Address())

	then.Bob().Operations().Last().
		WithTxID(internalTxID).
		WithTxStatus("BROADCASTED").
		WithValue(1).
		WithType("incoming").
		WithCounterparty(given.Alice().DefaultPaymail().Address())
}
