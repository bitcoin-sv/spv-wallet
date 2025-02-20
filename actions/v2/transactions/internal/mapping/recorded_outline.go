package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

// RecordedOutline maps domain RecordedOutline to api.ModelsRecordedOutline.
func RecordedOutline(r *txmodels.RecordedOutline) api.ModelsRecordedOutline {
	return api.ModelsRecordedOutline{
		TxID: r.TxID,
	}
}
