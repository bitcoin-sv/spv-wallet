package filter

// Page is a struct for handling paging parameters for search requests
type Page struct {
	Number int    `json:"page,omitempty"`
	Size   int    `json:"size,omitempty"`
	Sort   string `json:"sort,omitempty"`
	SortBy string `json:"sortBy,omitempty"`
}

// SearchParams is a generic struct for handling request parameters for search requests
type SearchParams[T any] struct {
	//nolint:staticcheck // SA5008 We want to reuse json tags also to mapstructure.
	Page Page `json:"paging,squash"`
	//nolint:staticcheck // SA5008 We want to reuse json tags also to mapstructure.
	Conditions T                      `json:"conditions,squash"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
