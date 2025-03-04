package integrationtests

import (
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

func TestHandlingARCCallback(t *testing.T) {
	tests := map[string]struct {
		txInfo         chainmodels.TXInfo
		expectStatus   txmodels.TxStatus
		beforeCallback func(t *testing.T, txInfo *chainmodels.TXInfo)
	}{
		"On Mined with BUMP": {
			txInfo: chainmodels.TXInfo{
				TXStatus:    chainmodels.Mined,
				BlockHeight: 885803,
				BlockHash:   "00000000000000000f0905597b6cac80031f0f56834e74dce1a714c682a9ed38",
				Timestamp:   time.Now(),
			},
			beforeCallback: calcBump,
			expectStatus:   txmodels.TxStatusMined,
		},
		"On SentToNetwork do nothing": {
			txInfo:       minimalTxInfo(chainmodels.SentToNetwork),
			expectStatus: txmodels.TxStatusBroadcasted,
		},
		"On SeenOnNetwork do nothing": {
			txInfo:       minimalTxInfo(chainmodels.SeenOnNetwork),
			expectStatus: txmodels.TxStatusBroadcasted,
		},
		"On Mined without BUMP don't change status": {
			txInfo: chainmodels.TXInfo{
				TXStatus:    chainmodels.Mined,
				BlockHash:   "00000000000000000f0905597b6cac80031f0f56834e74dce1a714c682a9ed38",
				BlockHeight: 885803,
				Timestamp:   time.Now(),
			},
			expectStatus: txmodels.TxStatusBroadcasted,
		},
		"On DoubleSpendAttempted mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.DoubleSpendAttempted),
			expectStatus: txmodels.TxStatusProblematic,
		},
		"On SeenInOrphanMempool mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.SeenInOrphanMempool),
			expectStatus: txmodels.TxStatusProblematic,
		},
		"On Rejected mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.Rejected),
			expectStatus: txmodels.TxStatusProblematic,
		},
		"On Unknown mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.Unknown),
			expectStatus: txmodels.TxStatusProblematic,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, when, then := testabilities.New(t)
			cleanup := given.StartedSPVWalletV2()
			defer cleanup()

			// when:
			receiveTxID := when.Alice().ReceivesFromExternal(10)

			// then:
			then.ARC().Broadcasted().
				WithTxID(receiveTxID).
				WithCallbackURL("https://example.com/transaction/broadcast/callback").
				WithCallbackToken(testengine.CallbackTestToken)

			// and:
			then.Alice().Operations().Last().
				WithTxID(receiveTxID).
				WithTxStatus("BROADCASTED")

			// given:
			test.txInfo.TxID = receiveTxID
			if test.beforeCallback != nil {
				test.beforeCallback(t, &test.txInfo)
			}

			// when:
			test.txInfo.TxID = receiveTxID
			if test.beforeCallback != nil {
				test.beforeCallback(t, &test.txInfo)
			}
			when.ARC().SendsCallback(test.txInfo)

			// then:
			then.Alice().Operations().Last().
				WithTxID(receiveTxID).
				WithTxStatus(string(test.expectStatus))

			// TODO: Assert for tx block height/hash and BEEF after Searching Transactions is implemented
		})
	}
}

func minimalTxInfo(txStatus chainmodels.TXStatus) chainmodels.TXInfo {
	return chainmodels.TXInfo{
		TXStatus:  txStatus,
		Timestamp: time.Now(),
	}
}

func calcBump(t *testing.T, txInfo *chainmodels.TXInfo) {
	t.Helper()

	txID, _ := chainhash.NewHashFromHex(txInfo.TxID)
	bump := trx.NewMerklePath(uint32(txInfo.BlockHeight), [][]*trx.PathElement{{
		&trx.PathElement{
			Hash:   txID,
			Offset: 0,
		},
	}})
	txInfo.MerklePath = bump.Hex()
}
