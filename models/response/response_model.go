package response

// PageDescription is a model that represents the page descriptor
type PageDescription struct {
	// Size is the number of elements on a single page
	Size int `json:"size"`
	// Number is the number of the page returned
	Number int `json:"number"`
	// TotalElements is the total number of elements in the returned collection
	TotalElements int `json:"totalElements"`
	// TotalPages is total number of pages returned
	TotalPages int `json:"totalPages"`
}

// PageModel is a model that represents the full JSON response
type PageModel[T any] struct {
	// Content is the collection of elements that serves as the content
	Content []*T `json:"content"`
	// Page is the page descriptor
	Page PageDescription `json:"page"`
}
