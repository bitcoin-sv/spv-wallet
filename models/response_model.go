package models

type PageDescription struct {
	Size          int `json:"size"`
	Number        int `json:"number"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

// TODO: do we need full hateos content data here?
type PageModel[T any] struct {
	Content []T             `json:"content"`
	Page    PageDescription `json:"page"`
}
