package integrationtests

import (
	"testing"
	"time"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
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
				Timestamp:   time.Now().Add(10 * time.Minute),
			},
			beforeCallback: calcBump,
			expectStatus:   txmodels.TxStatusMined,
		},
		"On SeenOnNetwork do nothing": {
			txInfo:       minimalTxInfo(chainmodels.SeenOnNetwork),
			expectStatus: txmodels.TxStatusBroadcasted,
		},
		"On Rejected mark as problematic": {
			txInfo:       minimalTxInfo(chainmodels.Rejected),
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

			if test.expectStatus == txmodels.TxStatusMined {
				then.Alice().Operations().Last().
					WithBlockHash(test.txInfo.BlockHash).
					WithBlockHeight(test.txInfo.BlockHeight)
			}

			// TODO: Assert for BEEF after Searching Transactions is implemented
		})
	}
}

func minimalTxInfo(txStatus chainmodels.TXStatus) chainmodels.TXInfo {
	return chainmodels.TXInfo{
		TXStatus:  txStatus,
		Timestamp: time.Now().Add(10 * time.Minute),
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
