package query_test

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestContextShouldGetQueryNestedMapSuccessfulParsing(t *testing.T) {
	var emptyQueryMap map[string]any
	veryDeepNesting := ""
	currentLv := make(map[string]any)
	veryDeepNestingResult := currentLv
	for i := 0; i < query.MaxNestedMapDepth; i++ {
		currKey := "nested" + strconv.Itoa(i)
		veryDeepNesting += "[" + currKey + "]"
		if i == query.MaxNestedMapDepth-1 {
			currentLv[currKey] = "value"
			continue
		}
		currentLv[currKey] = make(map[string]any)
		currentLv = currentLv[currKey].(map[string]any)
	}

	tests := map[string]struct {
		url            string
		expectedResult map[string]any
	}{
		"no query params": {
			url:            "",
			expectedResult: emptyQueryMap,
		},
		"single query param": {
			url: "?foo=bar",
			expectedResult: map[string]any{
				"foo": "bar",
			},
		},
		"empty key and value": {
			url: "?=",
			expectedResult: map[string]any{
				"": "",
			},
		},
		"empty key with some value value": {
			url: "?=value",
			expectedResult: map[string]any{
				"": "value",
			},
		},
		"single key with empty value": {
			url: "?key=",
			expectedResult: map[string]any{
				"key": "",
			},
		},
		"only keys": {
			url: "?foo&bar",
			expectedResult: map[string]any{
				"foo": "",
				"bar": "",
			},
		},
		"encoded & sign in value": {
			url: "?foo=bar%26baz",
			expectedResult: map[string]any{
				"foo": "bar&baz",
			},
		},
		"encoded = sign in value": {
			url: "?foo=bar%3Dbaz",
			expectedResult: map[string]any{
				"foo": "bar=baz",
			},
		},
		"multiple query param": {
			url: "?foo=bar&mapkey=value1",
			expectedResult: map[string]any{
				"foo":    "bar",
				"mapkey": "value1",
			},
		},
		"map query param": {
			url: "?mapkey[key]=value",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key": "value",
				},
			},
		},
		"multiple different value types in map query param": {
			url: "?mapkey[key1]=value1&mapkey[key2]=1&mapkey[key3]=true",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key1": "value1",
					"key2": "1",
					"key3": "true",
				},
			},
		},
		"multiple different value types in array value of map query param": {
			url: "?mapkey[key][]=value1&mapkey[key][]=1&mapkey[key][]=true",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key": []string{"value1", "1", "true"},
				},
			},
		},
		"nested map query param": {
			url: "?mapkey[key][nested][moreNested]=value",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key": map[string]any{
						"nested": map[string]any{
							"moreNested": "value",
						},
					},
				},
			},
		},
		"very deep nested map query param": {
			url: "?mapkey" + veryDeepNesting + "=value",
			expectedResult: map[string]any{
				"mapkey": veryDeepNestingResult,
			},
		},
		"map query param with explicit arrays accessors ([]) at the value level will return array": {
			url: "?mapkey[key][]=value1&mapkey[key][]=value2",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key": []string{"value1", "value2"},
				},
			},
		},
		"map query param with implicit arrays (duplicated key) at the value level will return only first value": {
			url: "?mapkey[key]=value1&mapkey[key]=value2",
			expectedResult: map[string]any{
				"mapkey": map[string]any{
					"key": "value1",
				},
			},
		},
		"array query param": {
			url: "?mapkey[]=value1&mapkey[]=value2",
			expectedResult: map[string]any{
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
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, dicts)
		})
	}
}

func TestContextShouldGetQueryNestedMapParsingError(t *testing.T) {
	tooDeepNesting := strings.Repeat("[nested]", query.MaxNestedMapDepth+1)

	tests := map[string]struct {
		url            string
		expectedResult map[string]any
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
		"too deep nesting of the map in query params": {
			url:   "?mapkey" + tooDeepNesting + "=value",
			error: "maximum depth [100] of nesting in map exceeded",
		},
		"setting value and nested map at the same level": {
			url:   "?mapkey[key]=value&mapkey[key][nested]=value1",
			error: "trying to set different types at the same key",
		},
		"setting array and nested map at the same level": {
			url:   "?mapkey[key][]=value&mapkey[key][nested]=value1",
			error: "trying to set different types at the same key",
		},
		"setting nested map and array at the same level": {
			url:   "?mapkey[key][nested]=value1&mapkey[key][]=value",
			error: "trying to set different types at the same key",
		},
		"setting array and value at the same level": {
			url:   "?key[]=value1&key=value2",
			error: "trying to set different types at the same key",
		},
		"setting value and array at the same level": {
			url:   "?key=value1&key[]=value2",
			error: "trying to set different types at the same key",
		},
		"setting array and nested map at same query param": {
			url:   "?mapkey[]=value1&mapkey[key]=value2",
			error: "trying to set different types at the same key",
		},
		"setting nested map and array at same query param": {
			url:   "?mapkey[key]=value2&mapkey[]=value1",
			error: "trying to set different types at the same key",
		},
		"] without [ in query param": {
			url:   "?mapkey]=value",
			error: "invalid access to map key",
		},
		"[ without ] in query param": {
			url:   "?mapkey[key=value",
			error: "invalid access to map key",
		},
		"[ without ] in query param with nested key": {
			url:   "?mapkey[key[nested]=value",
			error: "invalid access to map key",
		},
		"[[key]] in query param": {
			url:   "?mapkey[[key]]=value",
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

			dicts, err := query.ShouldGetQueryNestedMap(c)
			require.ErrorContains(t, err, test.error)
			require.Nil(t, dicts)
		})
	}
}

func TestContextShouldGetQueryNestedForKeySuccessfulParsing(t *testing.T) {
	var emptyQueryMap map[string]any

	tests := map[string]struct {
		url            string
		key            string
		expectedResult map[string]any
	}{
		"no searched map key in query string": {
			url:            "?foo=bar",
			key:            "mapkey",
			expectedResult: emptyQueryMap,
		},
		"searched map key after other query params": {
			url: "?foo=bar&mapkey[key]=value",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": "value",
			},
		},
		"searched map key before other query params": {
			url: "?mapkey[key]=value&foo=bar",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": "value",
			},
		},
		"single key in searched map key": {
			url: "?mapkey[key]=value",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": "value",
			},
		},
		"multiple keys in searched map key": {
			url: "?mapkey[key1]=value1&mapkey[key2]=value2&mapkey[key3]=value3",
			key: "mapkey",
			expectedResult: map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		"nested key in searched map key": {
			url: "?mapkey[foo][nested]=value1",
			key: "mapkey",
			expectedResult: map[string]any{
				"foo": map[string]any{
					"nested": "value1",
				},
			},
		},
		"multiple nested keys in single key of searched map key": {
			url: "?mapkey[foo][nested1]=value1&mapkey[foo][nested2]=value2",
			key: "mapkey",
			expectedResult: map[string]any{
				"foo": map[string]any{
					"nested1": "value1",
					"nested2": "value2",
				},
			},
		},
		"multiple keys with nested keys of searched map key": {
			url: "?mapkey[key1][nested]=value1&mapkey[key2][nested]=value2",
			key: "mapkey",
			expectedResult: map[string]any{
				"key1": map[string]any{
					"nested": "value1",
				},
				"key2": map[string]any{
					"nested": "value2",
				},
			},
		},
		"multiple levels of nesting in searched map key": {
			url: "?mapkey[key][nested][moreNested]=value1",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": map[string]any{
					"nested": map[string]any{
						"moreNested": "value1",
					},
				},
			},
		},
		"query keys similar to searched map key": {
			url: "?mapkey[key]=value&mapkeys[key1]=value1&mapkey1=foo",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": "value",
			},
		},
		"handle explicit arrays accessors ([]) at the value level": {
			url: "?mapkey[key][]=value1&mapkey[key][]=value2",
			key: "mapkey",
			expectedResult: map[string]any{
				"key": []string{"value1", "value2"},
			},
		},
		"implicit arrays (duplicated key) at the value level will return only first value": {
			url: "?mapkey[key]=value1&mapkey[key]=value2",
			key: "mapkey",
			expectedResult: map[string]any{
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
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, dicts)
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
			require.ErrorContains(t, err, test.error)
			require.Nil(t, dicts)
		})
	}
}
