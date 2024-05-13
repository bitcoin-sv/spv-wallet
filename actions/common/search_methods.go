package common

import (
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"math"
)

// Count object to use when returning a count of database query results
type Count struct {
	Content any `json:"content,omitempty"`
	Page
}

// Page object to use when limiting and sorting database query results
type Page struct {
	Size          int `json:"size,omitempty"`
	Number        int `json:"number,omitempty"`
	TotalElements int `json:"totalElements,omitempty"`
	TotalPages    int `json:"totalPages,omitempty"`
}

// LoadDefaultQueryParams will load the default query parameters
func LoadDefaultQueryParams() *datastore.QueryParams {
	return &datastore.QueryParams{
		Page:     1,
		PageSize: 10,
	}
}

// WrapCountResponse will wrap the content with the count and query parameters
func WrapCountResponse(content any, count int64, queryParams *datastore.QueryParams) Count {
	return Count{
		Content: content,
		Page: Page{
			Size:          queryParams.PageSize,
			Number:        queryParams.Page,
			TotalElements: int(count),
			TotalPages:    int(math.Ceil(float64(count) / float64(queryParams.PageSize))),
		},
	}
}
