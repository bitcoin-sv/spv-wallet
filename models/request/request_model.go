package request

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

const (
	DefaultPage      = 0
	DefaultSize      = 100
	DefaultSortOrder = "asc"
)

type SearchParams[T any] struct {
	Paging     Paging   `form:"paging"`
	Conditions T        `form:"conditions"`
	Metadata   Metadata `form:"metadata"`
}

type Paging struct {
	Page   *int    `form:"page,default=1"`
	Size   *int    `form:"size,default=10"`
	Sort   *string `form:"sort,default=asc"`
	SortBy *string `form:"sortBy,default=id"`
}

type Metadata map[string]string

// MapToFilterQueryParams - for compatibility
func (paging *Paging) MapToFilterQueryParams() *filter.QueryParams {
	return &filter.QueryParams{
		Page:          *paging.Page,
		PageSize:      *paging.Size,
		OrderByField:  *paging.SortBy,
		SortDirection: *paging.Sort,
	}
}

// UnStringify - give this a better name
func (m *Metadata) UnStringify() map[string]interface{} {
	m2 := make(map[string]interface{}, len(*m))
	for k, v := range *m {
		m2[k] = v
	}
	return m2
}

type Pageable struct {
	Page int
	Size int
	Sort Sort
}

type Order struct {
	Property  string
	Direction string // TODO: convert to enumerable
}

type Sort struct {
	Orders []Order
}

// func (data *QueryParamsData) AddPageDescription(c *gin.Context) *QueryParamsData {
// 	pageFromQueryParam, err := strconv.Atoi(c.Query("page"))
// 	fmt.Println("Page from request", pageFromQueryParam)
// 	if err != nil {
// 		pageFromQueryParam = DefaultPage
// 	}

// 	sizeFromQueryParam, err := strconv.Atoi(c.Query("size"))
// 	fmt.Println("Size from request", sizeFromQueryParam)
// 	if err != nil {
// 		sizeFromQueryParam = DefaultSize
// 	}

// 	sort := c.QueryArray("sort")
// 	fmt.Println("Sort from request", sort)
// 	data.Page = Pageable{
// 		Page: pageFromQueryParam,
// 		Size: sizeFromQueryParam,
// 		Sort: *createSortFromQueryParam(sort),
// 	}

// 	return data
// }

// TODO: handle default sort order
// TODO: what to do if sort parameter is broken?
func createSortFromQueryParam(sort []string) *Sort {
	orders := make([]Order, 0)
	const indexOfProperty = 0
	const indexOfDirection = 1
	for _, s := range sort {
		tokens := strings.Split(s, ",")
		if len(tokens) == 2 {
			orders = append(orders, Order{
				Property:  tokens[indexOfProperty],
				Direction: tokens[indexOfDirection],
			})
		} else if len(tokens) == 1 {
			orders = append(orders, Order{
				Property:  tokens[indexOfProperty],
				Direction: DefaultSortOrder,
			})
		}
	}

	return &Sort{Orders: orders}
}
