package models

// PagedResponse object to use when returning database records in paged format
type PagedResponse[content any] struct {
	// List of records for the response
	Content []content `json:"content"`
	// Pagination details
	Page Page `json:"page"`
}

// Page object to use when limiting and sorting database query results
type Page struct {
	// Field by which to order the results
	OrderByField *string `json:"orderByField"`
	// Direction in which to order the results ASC/DSC
	SortDirection *string `json:"sortDirection"`
	// Total count of elements
	TotalElements int64 `json:"totalElements"`
	// Total number of possible pages
	TotalPages int `json:"totalPages"`
	// Size of the page
	Size int `json:"size"`
	// Page number
	Number int `json:"number"`
}

// SearchContactsResponse is a response model for searching contacts
type SearchContactsResponse = PagedResponse[*Contact]

// ExclusiveStartKeyPage represents a paginated response for database records using Exclusive Start Key paging
type ExclusiveStartKeyPage[T any] struct {
	// List of records for the response
	Content T
	// Pagination details
	Page ExclusiveStartKeyPageInfo
}

// ExclusiveStartKeyPageInfo represents the pagination information for limiting and sorting database query results
type ExclusiveStartKeyPageInfo struct {
	// Field by which to order the results
	OrderByField *string `json:"orderByField,omitempty"` // Optional ordering field
	// Direction in which to order the results (ASC or DESC)
	SortDirection *string `json:"sortDirection,omitempty"` // Optional sort direction
	// Total count of elements
	TotalElements int `json:"totalElements"`
	// Size of the page or returned data
	Size int `json:"size"`
	// Last evaluated key returned from the database
	LastEvaluatedKey string `json:"lastEvaluatedKey"`
}
