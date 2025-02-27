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
		"On Mined with BUMP": {
			txInfo: chainmodels.TXInfo{
				TXStatus:    chainmodels.Mined,
				BlockHash:   "00000000000000000f0905597b6cac80031f0f56834e74dce1a714c682a9ed38",
				BlockHeight: 885803,
				Timestamp:   time.Now(),
				MerklePath:  "fe2b840d000902fd95010026e9ab42eb62d82f23a8940aeb02d58d69e92a6aa7d1f5b48674efb50707fc6dfd940102cd12118ddd85e4383901e13a04cd1179cec4740095eb70d2795f285f05aac5bd01cb007a63e80a3d3a937485e84f3f1a81670b626901cd705952d90d251a1c5ca907a1016400d463d7cc431fde3583bb4ba22089a27ac7cda8c5b5353ead60ee82a4b9e38c8d013300d9130482034bd2baeae736db5faeef985df314c91b164ccc36322c9c854061fe011800a15436aa01ce6e5f4290b974438700cfa689410ac9247bc04caee82f4a407032010d0082efd108a321a156023ad70c5601921df4b88fdd5402db21eb123870282458dd010700bea1a97aebcd5ff3b118fc64c85e2643c09426d22c00c700f84b8064b6c5745301020035ffd9b671117fcba0c027b0e4b966c68cb0c10994f241a5a647db51c2d0bdaa0100009f212e79950c6b827f6139d1c2a21bdcb6a275a0740c94b7e24795c4d4d7280a",
			},
			expectStatus: "MINED",
		},
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
