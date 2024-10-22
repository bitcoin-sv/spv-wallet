package common

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// GetPageFromQueryParams will return a Page object from the query parameters and count value
func GetPageFromQueryParams(queryParams *filter.QueryParams, count int64) models.Page {
	totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))
	page := models.Page{
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

// GetPageDescriptionFromSearchParams - returns a PageDescription based on the provided SearchParams
func GetPageDescriptionFromSearchParams(queryParams *datastore.QueryParams, count int64) response.PageDescription {
	totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))

	pageDescription := response.PageDescription{
		Size:          queryParams.PageSize,
		Number:        queryParams.Page,
		TotalElements: int(count),
		TotalPages:    totalPages,
	}

	return pageDescription
}

// MapToTypeContracts is a generic function that maps elements from one slice to another.
func MapToTypeContracts[T any, U any](input []T, mapper func(T) U) []U {
	output := make([]U, 0, len(input))
	for _, item := range input {
		output = append(output, mapper(item))
	}
	return output
}
