package filter

type Page struct {
	Number int    `json:"page,omitempty"`
	Size   int    `json:"size,omitempty"`
	Order  string `json:"order,omitempty"`
	SortBy string `json:"sortBy,omitempty"`
}

type SearchParams[T any] struct {
	Page       Page                   `json:"paging,squash"`
	Conditions T                      `json:"conditions,squash"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
