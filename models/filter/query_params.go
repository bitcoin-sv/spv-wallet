package filter

// QueryParams object to use when limiting and sorting database query results
type QueryParams struct {
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"page_size,omitempty"`
	OrderByField  string `json:"order_by_field,omitempty"`
	SortDirection string `json:"sort_direction,omitempty"`
}

func (qp *QueryParams) DefaultIfNil() {
	if qp == nil {
		qp = &QueryParams{
			Page:     1,
			PageSize: 10,
		}
	}
}

// DefaultQueryParams will return the default query parameters
func DefaultQueryParams() *QueryParams {
	return &QueryParams{
		Page:     1,
		PageSize: 10,
	}
}
