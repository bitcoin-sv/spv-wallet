package handlers

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// UserHandler is the handler for the user
type UserHandler = func(c *gin.Context, userContext *reqctx.UserContext)

// AsUser wraps the handler with the user context
func AsUser(handler UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := reqctx.GetUserContext(c)
		if userContext.IsAdmin() {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrAdminAuthOnUserEndpoint, nil)
			return
		}
		handler(c, userContext)
	}
}

// AdminHandler is the handler for admin's requests
type AdminHandler = func(c *gin.Context, _ *reqctx.AdminContext)

// AsAdmin wraps the handler with the AdminContext
func AsAdmin(handler AdminHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
			// The CheckSignature is based on the request body, so we can't use it for GET and other (no-body) requests
			if err := middleware.CheckSignature(c); err != nil {
				spverrors.AbortWithErrorResponse(c, err, nil)
				return
			}
		}
		userContext := reqctx.GetUserContext(c)
		if !userContext.IsAdmin() {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrNotAnAdminKey, nil)
			return
		}
		handler(c, reqctx.NewAdminContext())
	}
}

// AsAdminOrUser allows both admin and user to access the handler
func AsAdminOrUser(handler UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := reqctx.GetUserContext(c)
		handler(c, userContext)
	}
}
