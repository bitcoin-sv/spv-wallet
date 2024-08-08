package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPage      = 0
	DefaultSize      = 100
	DefaultSortOrder = "asc"
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
func GetPageDescriptionFromQueryParams(pageable *request.Pageable, count int64) response.PageDescription {
	totalPages := int(math.Ceil(float64(count) / float64(pageable.Size)))

	page := response.PageDescription{
		Size:          pageable.Size,
		Number:        pageable.Page,
		TotalElements: int(count),
		TotalPages:    totalPages,
	}

	return page
}

// ExtractPageableFromRequest
func ExtractPageableFromRequest(c *gin.Context) *request.Pageable {
	pageFromQueryParam, err := strconv.Atoi(c.Query("page"))
	fmt.Println("Page from request", pageFromQueryParam)
	if err != nil {
		pageFromQueryParam = DefaultPage
	}

	sizeFromQueryParam, err := strconv.Atoi(c.Query("size"))
	fmt.Println("Size from request", sizeFromQueryParam)
	if err != nil {
		sizeFromQueryParam = DefaultSize
	}

	sort := c.QueryArray("sort")
	fmt.Println("Sort from request", sort)
	return &request.Pageable{
		Page: pageFromQueryParam,
		Size: sizeFromQueryParam,
		Sort: *createSortFromQueryParam(sort),
	}
}

// TODO: handle default sort order
// TODO: what to do if sort parameter is broken?
func createSortFromQueryParam(sort []string) *request.Sort {
	orders := make([]request.Order, 0)
	const indexOfProperty = 0
	const indexOfDirection = 1
	for _, s := range sort {
		tokens := strings.Split(s, ",")
		if len(tokens) == 2 {
			orders = append(orders, request.Order{
				Property:  tokens[indexOfProperty],
				Direction: tokens[indexOfDirection],
			})
		} else if len(tokens) == 1 {
			orders = append(orders, request.Order{
				Property:  tokens[indexOfProperty],
				Direction: DefaultSortOrder,
			})
		}
	}

	return &request.Sort{Orders: orders}
}
