package common

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

// Count object to use when returning a count of database query results
type Count struct {
	Content any  `json:"content,omitempty"`
	Page    Page `json:"page,omitempty"`
}

// Page object to use when limiting and sorting database query results
type Page struct {
	Size          int    `json:"size"`
	Number        int    `json:"number"`
	TotalElements *int64 `json:"totalElements"`
	TotalPages    *int   `json:"totalPages"`
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
	totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))
	return Count{
		Content: content,
		Page: Page{
			Size:          queryParams.PageSize,
			Number:        queryParams.Page,
			TotalElements: &count,
			TotalPages:    &totalPages,
		},
	}
}

// WrapBasicSearchResponse will wrap the content without the count and query parameters
func WrapBasicSearchResponse(content any, size int) Count {
	return Count{
		Content: content,
		Page: Page{
			Size:   size,
			Number: 1,
		},
	}
}
