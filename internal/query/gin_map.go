package query

import "github.com/gin-gonic/gin"

//revive:disable:exported We want to mimic the gin API.

// ShouldGetQueryNestedMap returns a map from query params.
// In contrast to QueryMap it handles nesting in query maps like key[foo][bar]=value.
func ShouldGetQueryNestedMap(c *gin.Context) (dict map[string]any, err error) {
	return ShouldGetQueryNestedMapForKey(c, "")
}

// ShouldGetQueryNestedMapForKey returns a map from query params for a given query key.
// In contrast to QueryMap it handles nesting in query maps like key[foo][bar]=value.
// Similar to ShouldGetQueryNestedMap but it returns only the map for the given key.
func ShouldGetQueryNestedMapForKey(c *gin.Context, key string) (dict map[string]any, err error) {
	q := c.Request.URL.Query()
	return GetMap(q, key)
}

//revive:enable:exported
