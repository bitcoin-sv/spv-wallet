package common

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

// ConditionsModel is a generic model for handling conditions with metadata
type ConditionsModel[TFilter any] struct {
	// Custom conditions used for filtering the search results. Every field within the object is optional.
	Conditions TFilter `json:"conditions"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata *engine.Metadata `json:"metadata,omitempty" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// SearchModel is a generic model for handling searching with filters and metadata
type SearchModel[TFilter any] struct {
	ConditionsModel[TFilter]

	// Pagination and sorting options to streamline data exploration and analysis
	QueryParams *datastore.QueryParams `json:"params,omitempty" swaggertype:"object,string" example:"page:1,page_size:10,order_by_field:created_at,order_by_direction:desc"`
}
