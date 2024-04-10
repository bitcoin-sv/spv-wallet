package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// SearchRequestDestinationParameters is a struct for handling request parameters for search requests
type SearchRequestDestinationParameters struct {
	// Custom conditions used for filtering the search results
	Conditions filter.DestinationFilter `json:"conditions"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata *engine.Metadata `json:"metadata,omitempty" swaggertype:"object,string" example:"key:value,key2:value2"`
	// Pagination and sorting options to streamline data exploration and analysis
	QueryParams *datastore.QueryParams `json:"params,omitempty" swaggertype:"object,string" example:"page:1,page_size:10,order_by_field:created_at,order_by_direction:desc"`
}

// CountRequestDestinationParameters is a struct for handling request parameters for count requests
type CountRequestDestinationParameters struct {
	// Custom conditions used for filtering the search results
	Conditions map[string]interface{} `json:"conditions"  swaggertype:"object,string" example:"testColumn:testValue"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}
