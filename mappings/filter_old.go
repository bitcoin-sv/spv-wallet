package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// MapToQueryParams converts filter.QueryParams from models to matching datastore.QueryParams
func MapToQueryParams(model *filter.QueryParams) *datastore.QueryParams {
	if model == nil {
		return nil
	}
	return &datastore.QueryParams{
		Page:          model.Page,
		PageSize:      model.PageSize,
		OrderByField:  model.OrderByField,
		SortDirection: model.SortDirection,
	}
}
