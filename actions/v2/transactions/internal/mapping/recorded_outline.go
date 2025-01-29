package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// RecordedOutline maps domain RecordedOutline to response.RecordedOutline.
func RecordedOutline(r *txmodels.RecordedOutline) response.RecordedOutline {
	return response.RecordedOutline{
		TxID: r.TxID,
	}
}
