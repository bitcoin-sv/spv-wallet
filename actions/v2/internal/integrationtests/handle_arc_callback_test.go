package integrationtests

import (
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
)

func Test(t *testing.T) {
	tests := map[string]struct {
		txInfo       chainmodels.TXInfo
		expectStatus string
	}{
		"On SentToNetwork do nothing": {
			txInfo:       minimalTxInfo(chainmodels.SentToNetwork),
			expectStatus: "BROADCASTED",
		},
		"On SeenOnNetwork do nothing": {
			txInfo:       minimalTxInfo(chainmodels.SeenOnNetwork),
			expectStatus: "BROADCASTED",
		},
		"On DoubleSpendAttempted mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.DoubleSpendAttempted),
			expectStatus: "PROBLEMATIC",
		},
		"On SeenInOrphanMempool mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.SeenInOrphanMempool),
			expectStatus: "PROBLEMATIC",
		},
		"On Rejected mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.Rejected),
			expectStatus: "PROBLEMATIC",
		},
		"On Unknown mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.Unknown),
			expectStatus: "PROBLEMATIC",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, when, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletV2(testengine.WithARCCallback("https://example.com", testabilities.ARCCallbackToken))
			defer cleanup()

			// when:
			receiveTxID := when.Alice().ReceivesFromExternal(10)

			// then:
			then.ARC().Broadcasted().
				WithTxID(receiveTxID).
				WithCallbackURL("https://example.com/arc/broadcast/callback").
				WithCallbackToken(testabilities.ARCCallbackToken)

			// and:
			then.Alice().Operations().Last().
				WithTxID(receiveTxID).
				WithTxStatus("BROADCASTED")

			// when:
			test.txInfo.TxID = receiveTxID
			when.ARC().Callbacks(test.txInfo)

			// then:
			then.Alice().Operations().Last().
				WithTxID(receiveTxID).
				WithTxStatus(test.expectStatus)

			// TODO: Assert for tx block height and BEEF
		})
	}
}

func minimalTxInfo(txStatus chainmodels.TXStatus) chainmodels.TXInfo {
	return chainmodels.TXInfo{
		TXStatus:  txStatus,
		Timestamp: time.Now(),
	}
}
