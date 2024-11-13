package query

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestParseSearchParamsSuccessfully(t *testing.T) {
	tests := map[string]struct {
		url            string
		expectedResult filter.SearchParams[ExampleConditionsForTests]
	}{
		"empty query": {
			url:            "",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{},
		},
		"query page": {
			url: "?page=2&size=200&sort=asc&sortBy=id",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{
				Page: filter.Page{
					Number: 2,
					Size:   200,
					Sort:   "asc",
					SortBy: "id",
				},
			},
		},
		"query conditions model filter": {
			url: "?includeDeleted=true&createdRange[from]=2021-01-01T00:00:00Z&createdRange[to]=2021-01-02T00:00:00Z&updatedRange[from]=2021-02-01T00:00:00Z&updatedRange[to]=2021-02-02T00:00:00Z",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{
				Conditions: ExampleConditionsForTests{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: ptr(true),
						CreatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
		},
		"query conditions without nested structs": {
			url: "?xBoolean=true&xString=some%20string&xInt=5",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{
				Conditions: ExampleConditionsForTests{
					XBoolean: ptr(true),
					XString:  ptr("some string"),
					XInt:     ptr(5),
				},
			},
		},
		"query conditions nested struct": {
			url: "?nested[isNested]=true&nested[name]=some%20name&nested[number]=10",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{
				Conditions: ExampleConditionsForTests{
					Nested: &ExampleNestedConditionsForTests{
						IsNested: ptr(true),
						Name:     ptr("some name"),
						Number:   ptr(10),
					},
				},
			},
		},
		"query metadata": {
			url: "?metadata[key]=value1&metadata[key2][nested]=value2&metadata[key3][nested][]=value3&metadata[key3][nested][]=value4",
			expectedResult: filter.SearchParams[ExampleConditionsForTests]{
				Metadata: map[string]interface{}{
					"key": "value1",
					"key2": map[string]interface{}{
						"nested": "value2",
					},
					"key3": map[string]interface{}{
						"nested": []string{"value3", "value4"},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			require.NoError(t, err)

			c := &gin.Context{
				Request: &http.Request{
					URL: u,
				},
			}

			params, err := ParseSearchParams[ExampleConditionsForTests](c)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedResult, *params)
		})
	}
}

func TestNestingInArrayErrorCase(t *testing.T) {
	u, err := url.Parse("?metadata[key1][][key2]=value1&metadata[key1][][key2]=value2")
	require.NoError(t, err)

	c := &gin.Context{
		Request: &http.Request{
			URL: u,
		},
	}

	_, err = ParseSearchParams[ExampleConditionsForTests](c)
	require.ErrorContains(t, err, "unsupported array-like access to map key")
}

type ExampleConditionsForTests struct {
	// ModelFilter is a struct for handling typical request parameters for search requests
	//nolint:staticcheck // SA5008 - We want to reuse json tags also to mapstructure.
	filter.ModelFilter `json:",inline,squash"`
	XBoolean           *bool                            `json:"xBoolean,omitempty"`
	XString            *string                          `json:"xString,omitempty"`
	XInt               *int                             `json:"xInt,omitempty"`
	Nested             *ExampleNestedConditionsForTests `json:"nested,omitempty"`
}
type ExampleNestedConditionsForTests struct {
	IsNested *bool   `json:"isNested,omitempty"`
	Name     *string `json:"name,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

func ptr[T any](value T) *T {
	return &value
}
