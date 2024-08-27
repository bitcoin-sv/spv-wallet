package filter

type Page struct {
	Number int    `json:"page,omitempty" swaggertype:"integer" example:"2"`
	Size   int    `json:"size,omitempty" swaggertype:"integer" example:"5"`
	Order  string `json:"order,omitempty" swaggertype:"string" example:"desc"`
	SortBy string `json:"sortBy,omitempty" swaggertype:"string" example:"id"`
}

type SearchParams[T any] struct {
	Page       Page                   `json:"paging,squash"`
	Conditions T                      `json:"conditions,squash"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" swaggertype:"object,string" example:"key:value,key2:value2"`
}
