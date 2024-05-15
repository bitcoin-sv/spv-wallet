package common

import (
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"math"
)

// LoadDefaultQueryParams will load the default query parameters
func LoadDefaultQueryParams() *datastore.QueryParams {
	return &datastore.QueryParams{
		Page:     1,
		PageSize: 10,
	}
}

// GetPageFromQueryParams will return a Page object from the query parameters and count value
func GetPageFromQueryParams(queryParams *datastore.QueryParams, count int64) Page {
	totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))
	page := Page{
		Size:          queryParams.PageSize,
		Number:        queryParams.Page,
		TotalElements: count,
		TotalPages:    totalPages,
	}
	if queryParams.OrderByField != "" {
		page.OrderByField = &queryParams.OrderByField
	}
	if queryParams.SortDirection != "" {
		page.SortDirection = &queryParams.SortDirection
	}
	return page
}
