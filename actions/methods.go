package actions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/gin-gonic/gin"
)

// StatusOK is a basic response which sets the status to 200
func StatusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}

// NotFound handles all 404 requests
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, dictionary.GetError(dictionary.ErrorRequestNotFound, c.Request.RequestURI))
}

// MethodNotAllowed handles all 405 requests
func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, dictionary.GetError(dictionary.ErrorMethodNotAllowed, c.Request.Method, c.Request.RequestURI))
}
