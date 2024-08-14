package query_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestContextQueryNestedMap(t *testing.T) {
	var emptyQueryMap map[string]interface{}

	tests := map[string]struct {
		url            string
		expectedResult map[string]interface{}
		exists         bool
	}{
		"no searched map key in query string": {
			url:            "?foo=bar",
			expectedResult: emptyQueryMap,
			exists:         false,
		},
		"searched map key is not a map": {
			url:            "?mapkey=value",
			expectedResult: emptyQueryMap,
			exists:         false,
		},
		"searched map key is array": {
			url:            "?mapkey[]=value1&mapkey[]=value2",
			expectedResult: emptyQueryMap,
			exists:         false,
		},
		"searched map key with invalid map access": {
			url:            "?mapkey[key]nested=value",
			expectedResult: emptyQueryMap,
			exists:         false,
		},
		"searched map key with valid and invalid map access": {
			url: "?mapkey[key]invalidNested=value&mapkey[key][nested]=value1",
			expectedResult: map[string]interface{}{
				"key": map[string]interface{}{
					"nested": "value1",
				},
			},
			exists: true,
		},
		"searched map key after other query params": {
			url: "?foo=bar&mapkey[key]=value",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
			exists: true,
		},
		"searched map key before other query params": {
			url: "?mapkey[key]=value&foo=bar",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
			exists: true,
		},
		"single key in searched map key": {
			url: "?mapkey[key]=value",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
			exists: true,
		},
		"multiple keys in searched map key": {
			url: "?mapkey[key1]=value1&mapkey[key2]=value2&mapkey[key3]=value3",
			expectedResult: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			exists: true,
		},
		"nested key in searched map key": {
			url: "?mapkey[foo][nested]=value1",
			expectedResult: map[string]interface{}{
				"foo": map[string]interface{}{
					"nested": "value1",
				},
			},
			exists: true,
		},
		"multiple nested keys in single key of searched map key": {
			url: "?mapkey[foo][nested1]=value1&mapkey[foo][nested2]=value2",
			expectedResult: map[string]interface{}{
				"foo": map[string]interface{}{
					"nested1": "value1",
					"nested2": "value2",
				},
			},
			exists: true,
		},
		"multiple keys with nested keys of searched map key": {
			url: "?mapkey[key1][nested]=value1&mapkey[key2][nested]=value2",
			expectedResult: map[string]interface{}{
				"key1": map[string]interface{}{
					"nested": "value1",
				},
				"key2": map[string]interface{}{
					"nested": "value2",
				},
			},
			exists: true,
		},
		"multiple levels of nesting in searched map key": {
			url: "?mapkey[key][nested][moreNested]=value1",
			expectedResult: map[string]interface{}{
				"key": map[string]interface{}{
					"nested": map[string]interface{}{
						"moreNested": "value1",
					},
				},
			},
			exists: true,
		},
		"query keys similar to searched map key": {
			url: "?mapkey[key]=value&mapkeys[key1]=value1&mapkey1=foo",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
			exists: true,
		},
	}
	for name, test := range tests {
		t.Run("getQueryMap: "+name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			require.NoError(t, err)

			c := &gin.Context{
				Request: &http.Request{
					URL: u,
				},
			}

			dicts, exists := query.GetQueryNestedMap(c, "mapkey")
			require.Equal(t, test.expectedResult, dicts)
			require.Equal(t, test.exists, exists)
		})
		t.Run("queryMap: "+name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			require.NoError(t, err)

			c := &gin.Context{
				Request: &http.Request{
					URL: u,
				},
			}

			dicts := query.QueryNestedMap(c, "mapkey")
			require.Equal(t, test.expectedResult, dicts)
		})
	}
}
