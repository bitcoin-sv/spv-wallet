package common

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/request"
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

// GetPageDescriptionFromQueryParams
func GetPageDescriptionFromQueryParams(qp *request.QueryParamsData, count int64) response.PageDescription {
	totalPages := int(math.Ceil(float64(count) / float64(qp.Page.Size)))

	page := response.PageDescription{
		Size:          qp.Page.Size,
		Number:        qp.Page.Page,
		TotalElements: int(count),
		TotalPages:    totalPages,
	}

	return page
}
