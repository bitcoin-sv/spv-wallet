package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

const (
	DefaultPage      = 0
	DefaultSize      = 100
	DefaultSortOrder = "asc"
)

type QueryParamsSchema struct {
	Conditions Conditions `schema:"conditions"`
	Metadata   Metadata   `schema:"metadata"`
	// Page       Pageable `schema:"conditions"`
	Page int  `schema:"page"`
	Size int  `schema:"size"`
	Sort Sort `schema:"sort"`
}

type QueryParamsData struct {
	Conditions Conditions
	Metadata   Metadata
	Page       Pageable
}

type Conditions map[string]interface{}

type Metadata map[string]string

type Pageable struct {
	Page int
	Size int
	Sort Sort
}

type Order struct {
	Property  string
	Direction string // TODO: convert to enumerable
}

type Sort struct {
	Orders []Order
}

// func ExtractQueryParamsFromRequest(c *gin.Context) url.Values {
// 	return c.Request.URL.Query()
// }

func ExtractQueryParamsFromRequest(r *http.Request) *QueryParamsSchema {
	qp := r.URL.Query()
	data := &QueryParamsSchema{}

	fmt.Printf("Query params from request: %+v\n", qp)

	decoder := schema.NewDecoder()
	err := decoder.Decode(data, qp)
	if err != nil {
		fmt.Println("Could not decode the url query params", err)
	}

	return data
}

// CreateRequestQueryParams
func CreateQueryParamsData() *QueryParamsData {
	data := QueryParamsData{}

	return &data
}

// AddConditions
func (data *QueryParamsData) AddConditions(qp url.Values) *QueryParamsData {
	fmt.Println("Conditions from request", qp)
	// qp.Conditions = Conditions["test":"test"]

	return data
}

// AddMetadata
func (data *QueryParamsData) AddMetadata(qp url.Values) *QueryParamsData {
	fmt.Println("Metadata from request", qp)
	// if err != nil {
	// 	// metadataFromQueryParam = map
	// }

	// qp.Metadata = Metadata["test":"test"]

	return data
}

// AddPageDescription
// func (data *QueryParamsData) AddPageDescription() *QueryParamsData {}

// func (data *QueryParamsData) AddPageDescription(c *gin.Context) *QueryParamsData {
// 	pageFromQueryParam, err := strconv.Atoi(c.Query("page"))
// 	fmt.Println("Page from request", pageFromQueryParam)
// 	if err != nil {
// 		pageFromQueryParam = DefaultPage
// 	}

// 	sizeFromQueryParam, err := strconv.Atoi(c.Query("size"))
// 	fmt.Println("Size from request", sizeFromQueryParam)
// 	if err != nil {
// 		sizeFromQueryParam = DefaultSize
// 	}

// 	sort := c.QueryArray("sort")
// 	fmt.Println("Sort from request", sort)
// 	data.Page = Pageable{
// 		Page: pageFromQueryParam,
// 		Size: sizeFromQueryParam,
// 		Sort: *createSortFromQueryParam(sort),
// 	}

// 	return data
// }

// TODO: handle default sort order
// TODO: what to do if sort parameter is broken?
func createSortFromQueryParam(sort []string) *Sort {
	orders := make([]Order, 0)
	const indexOfProperty = 0
	const indexOfDirection = 1
	for _, s := range sort {
		tokens := strings.Split(s, ",")
		if len(tokens) == 2 {
			orders = append(orders, Order{
				Property:  tokens[indexOfProperty],
				Direction: tokens[indexOfDirection],
			})
		} else if len(tokens) == 1 {
			orders = append(orders, Order{
				Property:  tokens[indexOfProperty],
				Direction: DefaultSortOrder,
			})
		}
	}

	return &Sort{Orders: orders}
}
