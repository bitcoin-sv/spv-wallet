package metrics

import (
	"github.com/gin-gonic/gin"
)

var notFoundContextKey = "routeNotFound"

func requestMetricsMiddleware() gin.HandlerFunc {
	if metrics, enabled := Get(); enabled {
		return func(c *gin.Context) {
			tracker := metrics.httpRequests.Track(c.Request.Method, c.Request.URL.Path)
			tracker.Start()
			defer func() {
				if _, noRoute := c.Get(notFoundContextKey); noRoute {
					tracker.EndWithNoRoute()
				} else {
					// note that the status code will be correct only if higher middleware doesn't change the status code;
					// order of middlewares matters
					tracker.End(c.Writer.Status())
				}
			}()

			c.Next()
		}
	}

	return func(c *gin.Context) {
		c.Next()
	}
}

// NoRoute is a middleware to set a flag in the context if the route is actually not found
// this is needed to distinguish no-route 404 from other 404s
func NoRoute(c *gin.Context) {
	c.Set(notFoundContextKey, true)
	c.Next()
}
