package handlers

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// UserHandler is the handler for the user
type UserHandler = func(c *gin.Context, userContext *reqctx.UserContext)

// UserHandlerWithAppConfig is the handler for the user that also takes appConfig
type UserHandlerWithAppConfig = func(c *gin.Context, appConfig *config.AppConfig)

// UserHandlerWithXPub is the handler for the user who has authorized using xPub
type UserHandlerWithXPub = func(c *gin.Context, userContext *reqctx.UserContext, xpub string)

// AsUser wraps the handler with the user context. User can be authorized by xPub or accessKey
func AsUser(handler UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := reqctx.GetUserContext(c)
		if userContext.GetAuthType() == reqctx.AuthTypeAdmin {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrAdminAuthOnUserEndpoint, nil)
			return
		}
		handler(c, userContext)
	}
}

// AsUserWithAppConfig wraps the handler with the user context. User can be authorized by xPub or accessKey
// it also passes appConfig to the handler
func AsUserWithAppConfig(handler UserHandlerWithAppConfig, appConfig *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := reqctx.GetUserContext(c)
		if userContext.GetAuthType() == reqctx.AuthTypeAdmin {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrAdminAuthOnUserEndpoint, nil)
			return
		}
		handler(c, appConfig)
	}
}

// AdminHandler is the handler for admin's requests
type AdminHandler = func(c *gin.Context, _ *reqctx.AdminContext)

// AsAdmin wraps the handler with the AdminContext
func AsAdmin(handler AdminHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := reqctx.GetUserContext(c)
		if userContext.GetAuthType() != reqctx.AuthTypeAdmin {
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
