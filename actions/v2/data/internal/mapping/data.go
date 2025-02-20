package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
)

// DataResponse maps a domain data model to a response model
func DataResponse(data *datamodels.Data) api.ModelsData {
	return api.ModelsData{
		Id:   data.ID(),
		Blob: string(data.Blob),
	}
}
