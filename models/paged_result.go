package models

// PagedResult is a generic struct for paginated results.
type PagedResult[T any] struct {
	Content         []*T
	PageDescription PageDescription
}

type PageDescription struct {
	Size          int
	Number        int
	TotalElements int
	TotalPages    int
}
