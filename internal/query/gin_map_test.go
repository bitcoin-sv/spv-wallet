package query_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestContextShouldGetQueryNestedMapSuccessfulParsing(t *testing.T) {
	var emptyQueryMap map[string]interface{}

	tests := map[string]struct {
		url            string
		expectedResult map[string]interface{}
	}{
		"no query params": {
			url:            "",
			expectedResult: emptyQueryMap,
		},
		"single query param": {
			url: "?foo=bar",
			expectedResult: map[string]interface{}{
				"foo": "bar",
			},
		},
		"multiple query param": {
			url: "?foo=bar&mapkey=value1",
			expectedResult: map[string]interface{}{
				"foo":    "bar",
				"mapkey": "value1",
			},
		},
		"map query param": {
			url: "?mapkey[key]=value",
			expectedResult: map[string]interface{}{
				"mapkey": map[string]interface{}{
					"key": "value",
				},
			},
		},
		"nested map query param": {
			url: "?mapkey[key][nested][moreNested]=value",
			expectedResult: map[string]interface{}{
				"mapkey": map[string]interface{}{
					"key": map[string]interface{}{
						"nested": map[string]interface{}{
							"moreNested": "value",
						},
					},
				},
			},
		},
		"map query param with explicit arrays accessors ([]) at the value level will return array": {
			url: "?mapkey[key][]=value1&mapkey[key][]=value2",
			expectedResult: map[string]interface{}{
				"mapkey": map[string]interface{}{
					"key": []string{"value1", "value2"},
				},
			},
		},
		"map query param with implicit arrays (duplicated key) at the value level will return only first value": {
			url: "?mapkey[key]=value1&mapkey[key]=value2",
			expectedResult: map[string]interface{}{
				"mapkey": map[string]interface{}{
					"key": "value1",
				},
			},
		},
		"array query param": {
			url: "?mapkey[]=value1&mapkey[]=value2",
			expectedResult: map[string]interface{}{
				"mapkey": []string{"value1", "value2"},
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

			dicts, err := query.ShouldGetQueryNestedMap(c)
			require.Equal(t, test.expectedResult, dicts)
			require.NoError(t, err)
		})
	}
}

func TestContextShouldGetQueryNestedMapParsingError(t *testing.T) {
	tests := map[string]struct {
		url            string
		expectedResult map[string]interface{}
		error          string
	}{
		"searched map key with invalid map access": {
			url:   "?mapkey[key]nested=value",
			error: "invalid access to map key",
		},
		"searched map key with array accessor in the middle": {
			url:   "?mapkey[key][][nested]=value",
			error: "unsupported array-like access to map key",
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

			dicts, err := query.ShouldGetQueryNestedMap(c)
			require.Nil(t, dicts)
			require.ErrorContains(t, err, test.error)
		})
	}
}

func TestContextShouldGetQueryNestedForKeySuccessfulParsing(t *testing.T) {
	var emptyQueryMap map[string]interface{}

	tests := map[string]struct {
		url            string
		key            string
		expectedResult map[string]interface{}
	}{
		"no searched map key in query string": {
			url:            "?foo=bar",
			key:            "mapkey",
			expectedResult: emptyQueryMap,
		},
		"searched map key after other query params": {
			url: "?foo=bar&mapkey[key]=value",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
		},
		"searched map key before other query params": {
			url: "?mapkey[key]=value&foo=bar",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
		},
		"single key in searched map key": {
			url: "?mapkey[key]=value",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
		},
		"multiple keys in searched map key": {
			url: "?mapkey[key1]=value1&mapkey[key2]=value2&mapkey[key3]=value3",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		"nested key in searched map key": {
			url: "?mapkey[foo][nested]=value1",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"foo": map[string]interface{}{
					"nested": "value1",
				},
			},
		},
		"multiple nested keys in single key of searched map key": {
			url: "?mapkey[foo][nested1]=value1&mapkey[foo][nested2]=value2",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"foo": map[string]interface{}{
					"nested1": "value1",
					"nested2": "value2",
				},
			},
		},
		"multiple keys with nested keys of searched map key": {
			url: "?mapkey[key1][nested]=value1&mapkey[key2][nested]=value2",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key1": map[string]interface{}{
					"nested": "value1",
				},
				"key2": map[string]interface{}{
					"nested": "value2",
				},
			},
		},
		"multiple levels of nesting in searched map key": {
			url: "?mapkey[key][nested][moreNested]=value1",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": map[string]interface{}{
					"nested": map[string]interface{}{
						"moreNested": "value1",
					},
				},
			},
		},
		"query keys similar to searched map key": {
			url: "?mapkey[key]=value&mapkeys[key1]=value1&mapkey1=foo",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": "value",
			},
		},
		"handle explicit arrays accessors ([]) at the value level": {
			url: "?mapkey[key][]=value1&mapkey[key][]=value2",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": []string{"value1", "value2"},
			},
		},
		"implicit arrays (duplicated key) at the value level will return only first value": {
			url: "?mapkey[key]=value1&mapkey[key]=value2",
			key: "mapkey",
			expectedResult: map[string]interface{}{
				"key": "value1",
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

			dicts, err := query.ShouldGetQueryNestedMapForKey(c, test.key)
			require.Equal(t, test.expectedResult, dicts)
			require.NoError(t, err)
		})
	}
}

func TestContextShouldGetQueryNestedForKeyParsingError(t *testing.T) {
	tests := map[string]struct {
		url   string
		key   string
		error string
	}{

		"searched map key is value not a map": {
			url:   "?mapkey=value",
			key:   "mapkey",
			error: "invalid access to map",
		},
		"searched map key is array": {
			url:   "?mapkey[]=value1&mapkey[]=value2",
			key:   "mapkey",
			error: "invalid access to map",
		},
		"searched map key with invalid map access": {
			url:   "?mapkey[key]nested=value",
			key:   "mapkey",
			error: "invalid access to map key",
		},
		"searched map key with valid and invalid map access": {
			url:   "?mapkey[key]invalidNested=value&mapkey[key][nested]=value1",
			key:   "mapkey",
			error: "invalid access to map key",
		},
		"searched map key with valid before invalid map access": {
			url:   "?mapkey[key][nested]=value1&mapkey[key]invalidNested=value",
			key:   "mapkey",
			error: "invalid access to map key",
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

			dicts, err := query.ShouldGetQueryNestedMapForKey(c, test.key)
			require.Nil(t, dicts)
			require.ErrorContains(t, err, test.error)
		})
	}
}
