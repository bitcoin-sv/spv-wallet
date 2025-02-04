package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// DataResponse maps a domain data model to a response model
func DataResponse(data *datamodels.Data) response.Data {
	return response.Data{
		ID:   data.ID(),
		Blob: string(data.Blob),
	}
}
