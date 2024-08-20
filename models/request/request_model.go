package request

import (
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
	Size   *int    `form:"size,default=100"`
	Sort   *string `form:"sort,default=asc"`
	SortBy *string `form:"sortBy,default=id"`
}

type Metadata map[string]interface{}

// MapToFilterQueryParams - for compatibility
func (paging *Paging) MapToFilterQueryParams() *filter.QueryParams {
	return &filter.QueryParams{
		Page:          *paging.Page,
		PageSize:      *paging.Size,
		OrderByField:  *paging.SortBy,
		SortDirection: *paging.Sort,
	}
}
