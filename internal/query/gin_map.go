package query

import "github.com/gin-gonic/gin"

// QueryNestedMap returns a map for a given query key.
// In contrast to QueryMap it handles nesting in query maps like key[foo][bar]=value.
func QueryNestedMap(c *gin.Context, key string) (dicts map[string]interface{}) {
	dicts, _ = GetQueryNestedMap(c, key)
	return
}

// GetQueryNestedMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
// In contrast to GetQueryMap it handles nesting in query maps like key[foo][bar]=value.
func GetQueryNestedMap(c *gin.Context, key string) (map[string]interface{}, bool) {
	q := c.Request.URL.Query()
	return GetMap(q, key)
}
