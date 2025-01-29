package common

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

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
