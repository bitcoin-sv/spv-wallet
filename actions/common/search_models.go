package common

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/models"
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

// PagedResponse object to use when returning database records in paged format
type PagedResponse[content any] struct {
	// List of records for the response
	Content []content `json:"content"`
	// Pagination details
	Page Page `json:"page"`
}

// Page object to use when limiting and sorting database query results
type Page struct {
	// Field by which to order the results
	OrderByField *string `json:"orderByField"`
	// Direction in which to order the results ASC/DSC
	SortDirection *string `json:"sortDirection"`
	// Total count of elements
	TotalElements int64 `json:"totalElements"`
	// Total number of possible pages
	TotalPages int `json:"totalPages"`
	// Size of the page
	Size int `json:"size"`
	// Page number
	Number int `json:"number"`
}

// SearchContactsResponse is a response model for searching contacts
type SearchContactsResponse = PagedResponse[*models.Contact]
