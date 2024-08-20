package common

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/gin-gonic/gin"
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
func GetPageDescriptionFromSearchParams[T any](params *request.SearchParams[T], count int64) response.PageDescription {
	totalPages := int(math.Ceil(float64(count) / float64(*params.Paging.Size)))

	page := response.PageDescription{
		Size:          *params.Paging.Size,
		Number:        *params.Paging.Page,
		TotalElements: int(count),
		TotalPages:    totalPages,
	}

	return page
}

// GetSearchParams - returns a SearchParams struct based on the QueryString
func GetSearchParams[T any](c *gin.Context, _ T) (request.SearchParams[T], error) {
	var queryParams request.SearchParams[T]

	// Bind the query parameters to the struct
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		return queryParams, spverrors.Wrapf(err, "Cannot bind query params")
	}

	queryParams.Metadata = query.QueryNestedMap(c, "metadata")
	return queryParams, nil
}
