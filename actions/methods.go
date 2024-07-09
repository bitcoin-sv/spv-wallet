package actions

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusOK is a basic response which sets the status to 200
func StatusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}

// NotFound handles all 404 requests
func NotFound(c *gin.Context) {
	spverrors.ErrorResponse(c, spverrors.ErrRouteNotFound, nil)
}

// MethodNotAllowed handles all 405 requests
func MethodNotAllowed(c *gin.Context) {
	spverrors.ErrorResponse(c, spverrors.ErrRouteMethodNotAllowed, nil)
}
