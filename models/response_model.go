package models

type PageDescription struct {
	Size          int `json:"size"`
	Number        int `json:"number"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

type PageModel[T any] struct {
	Content []*T            `json:"content"`
	Page    PageDescription `json:"page"`
}
