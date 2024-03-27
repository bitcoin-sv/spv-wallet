package datastore

import (
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
)

// QueryParams object to use when limiting and sorting database query results
type QueryParams struct {
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"page_size,omitempty"`
	OrderByField  string `json:"order_by_field,omitempty"`
	SortDirection string `json:"sort_direction,omitempty"`
}

// MarshalQueryParams will marshal the custom type
func MarshalQueryParams(m QueryParams) graphql.Marshaler {
	if m.Page == 0 && m.PageSize == 0 && m.OrderByField == "" && m.SortDirection == "" {
		return graphql.Null
	}
	return graphql.MarshalAny(m)
}

// UnmarshalQueryParams will unmarshal the custom type
func UnmarshalQueryParams(v interface{}) (QueryParams, error) {
	if v == nil {
		return QueryParams{}, nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return QueryParams{}, err
	}

	var q QueryParams
	if err = json.Unmarshal(data, &q); err != nil {
		return QueryParams{}, err
	}

	return q, nil
}
