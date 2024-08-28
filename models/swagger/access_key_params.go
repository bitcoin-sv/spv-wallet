package swagger

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type CommonFilteringQueryParams struct {
	filter.Page `json:",inline"`

	CreatedRangeFrom *time.Time `json:"createdRange[from],omitempty" example:"2024-02-26T11:01:28Z"`
	CreatedRangeTo   *time.Time `json:"createdRange[to],omitempty" example:"2024-02-26T11:01:28Z"`

	UpdatedRangeFrom *time.Time `json:"updatedRange[from],omitempty" example:"2024-02-26T11:01:28Z"`
	UpdatedRangeTo   *time.Time `json:"updatedRange[to],omitempty" example:"2024-02-26T11:01:28Z"`

	// Metadata is a list of key-value pairs that can be used to filter the results. !ATTENTION! Unfortunately this parameter won't work from swagger UI.
	Metadata []string `json:"metadata,omitempty" swaggertype:"array,string" example:"metadata[key]=value,metadata[key2]=value2"`
}
