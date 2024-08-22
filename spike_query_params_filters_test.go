package spvwallet_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const expected = `{"Conditions":{"includeDeleted":false,"createdRange":{"from":"2001-02-03T04:05:06Z","to":"2002-03-04T05:06:07Z"},"includeOthers":false,"executedRange":{"from":"2003-04-05T06:07:08Z","to":"2004-05-06T07:08:09Z"}}}`

// ========= some example query params structs ======

type Paging struct {
	Page   *int    `form:"page,default=1"`
	Size   *int    `form:"size,default=10"`
	Sort   *string `form:"sort,default=asc"`
	SortBy *string `form:"sortBy,default=id"`
}

type SearchParams[T any] struct {
	// Paging     Paging                 `form:"paging"`
	Conditions T `form:"conditions"`
	// Metadata   map[string]interface{} `form:"-"`
}

// ========= EXAMPLE JSON solution ======

type JsonRangeConditions struct {
	JsonModelFilter

	IncludeOthers *bool `form:"includeOthers,default=false" json:"includeOthers,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	ExecutedRange *filter.TimeRange `form:"executedRange" json:"executedRange,omitempty"`
}

// Simplified ModelFilter struct
type JsonModelFilter struct {
	// IncludeDeleted is a flag whether or not to include deleted items in the search results
	IncludeDeleted *bool `form:"includeDeleted,default=false" json:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	// CreatedRange specifies the time range when a record was created.
	CreatedRange *filter.TimeRange `form:"createdRange" json:"createdRange,omitempty"`
}

func TestJsonRangeConditions(t *testing.T) {
	// given
	ginCtx := makeRequest(`?createdRange={"from":"2001-02-03T04:05:06Z","to":"2002-03-04T05:06:07Z"}&executedRange={"from":"2003-04-05T06:07:08Z","to":"2004-05-06T07:08:09Z"}`)

	// when
	var queryParams SearchParams[JsonRangeConditions]
	err := ginCtx.ShouldBindQuery(&queryParams)
	if err != nil {
		panic(err)
	}

	// then
	res, err := json.Marshal(queryParams)
	require.Equal(t, expected, string(res))
}

// === EXAMPLE multistruct solution ===

type MultiStructRangeConditions struct {
	MultiStructModelFilter

	IncludeOthers *bool `form:"includeOthers,default=false" json:"includeOthers,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	ExecutedRange *ExecutedTimeRange `form:"executedRange" json:"executedRange,omitempty"`
}

type MultiStructModelFilter struct {
	// IncludeDeleted is a flag whether or not to include deleted items in the search results
	IncludeDeleted *bool `form:"includeDeleted,default=false" json:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	// CreatedRange specifies the time range when a record was created.
	CreatedRange *CreatedTimeRange `form:"createdRange" json:"createdRange,omitempty"`
}

type CreatedTimeRange struct {
	// From specifies the start time of the range. It's optional and can be nil.
	From *time.Time `form:"createdRange[from]" json:"from,omitempty" example:"2024-02-26T11:01:28Z"`
	// To specifies the end time of the range. It's optional and can be nil.
	To *time.Time `form:"createdRange[to]" json:"to,omitempty" example:"2024-02-26T11:01:28Z"`
}

type ExecutedTimeRange struct {
	// From specifies the start time of the range. It's optional and can be nil.
	From *time.Time `form:"executedRange[from]" json:"from,omitempty" example:"2024-02-26T11:01:28Z"`
	// To specifies the end time of the range. It's optional and can be nil.
	To *time.Time `form:"executedRange[to]" json:"to,omitempty" example:"2024-02-26T11:01:28Z"`
}

func TestMultiStructRangeConditions(t *testing.T) {
	// given
	ginCtx := makeRequest(`?createdRange[from]=2001-02-03T04:05:06Z&createdRange[to]=2002-03-04T05:06:07Z&executedRange[from]=2003-04-05T06:07:08Z&executedRange[to]=2004-05-06T07:08:09Z`)

	// when
	var queryParams SearchParams[MultiStructRangeConditions]
	err := ginCtx.ShouldBindQuery(&queryParams)
	if err != nil {
		panic(err)
	}

	// then
	res, err := json.Marshal(queryParams)
	require.Equal(t, expected, string(res))
}

// === EXAMPLE custom binder solution ===

type CustomBinderRangeConditions struct {
	CustomBinderModelFilter

	IncludeOthers *bool `form:"includeOthers,default=false" json:"includeOthers,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	ExecutedRange *CustomBinderTimeRange `form:"executedRange" json:"executedRange,omitempty"`
}

// Simplified ModelFilter struct
type CustomBinderModelFilter struct {
	// IncludeDeleted is a flag whether or not to include deleted items in the search results
	IncludeDeleted *bool `form:"includeDeleted,default=false" json:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	// CreatedRange specifies the time range when a record was created.
	CreatedRange *CustomBinderTimeRange `form:"createdRange" json:"createdRange,omitempty"`
}

type CustomBinderTimeRange struct {
	// From specifies the start time of the range. It's optional and can be nil.
	From *time.Time `form:"executedRange[from]" json:"from,omitempty" example:"2024-02-26T11:01:28Z"`
	// To specifies the end time of the range. It's optional and can be nil.
	To *time.Time `form:"executedRange[to]" json:"to,omitempty" example:"2024-02-26T11:01:28Z"`
}

func (t *CustomBinderTimeRange) UnmarshalParam(param string) error {
	// given argument "param" value = "from:2024-02-26T11:01:28Z,to:2024-02-26T11:01:28Z"
	// split by comma, then by colon and assign to From and To to t struct
	params := strings.Split(param, ",")
	for _, p := range params {
		kv := strings.Split(p, "~")
		val, err := time.Parse(time.RFC3339, kv[1])
		if err != nil {
			return err
		}
		switch kv[0] {
		case "from":
			t.From = &val
		case "to":
			t.To = &val
		default:
			return errors.New("invalid value")
		}
	}
	return nil
}

func TestCustomBinderRangeConditions(t *testing.T) {
	// given
	ginCtx := makeRequest(`?createdRange=from~2001-02-03T04:05:06Z,to~2002-03-04T05:06:07Z&executedRange=from~2003-04-05T06:07:08Z,to~2004-05-06T07:08:09Z`)

	// when
	var queryParams SearchParams[CustomBinderRangeConditions]
	err := ginCtx.ShouldBindQuery(&queryParams)
	if err != nil {
		panic(err)
	}

	// then
	res, err := json.Marshal(queryParams)
	require.Equal(t, expected, string(res))
}

//  ====== test helpers ======

func makeRequest(query string) *gin.Context {
	uri, err := url.Parse(query)
	if err != nil {
		panic(err)
	}
	ginCtx := gin.Context{
		Request: &http.Request{
			URL: uri,
		},
	}
	return &ginCtx
}
