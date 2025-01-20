package api

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

var securedMiddlewares = []MiddlewareFunc{
	(MiddlewareFunc)(middleware.AuthMiddleware()),
	(MiddlewareFunc)(middleware.CheckSignatureMiddleware()),
}

// SignatureAuthWithScopes checks for scopes and runs auth&signature middlewares
func SignatureAuthWithScopes() MiddlewareFunc {
	return func(c *gin.Context) {
		if scope, ok := c.Get(SignatureAuthScopes); ok {
			var userType string
			if scopes := scope.([]string); len(scopes) == 0 {
				spverrors.ErrorResponse(c, spverrors.ErrMissingAuthHeader, reqctx.Logger(c))
				return
			} else {
				userType = scopes[0]
			}

			if userType != "admin" && userType != "user" {
				spverrors.ErrorResponse(c, spverrors.ErrAuthorization, reqctx.Logger(c))
				return
			}

			for _, mid := range securedMiddlewares {
				mid(c)
				if c.IsAborted() {
					return
				}
			}
			userCtx := reqctx.GetUserContext(c)
			if userType == "admin" && userCtx.AuthType != reqctx.AuthTypeAdmin {
				spverrors.ErrorResponse(c, spverrors.ErrNotAnAdminKey, reqctx.Logger(c))
				return
			} else if userType == "user" && userCtx.AuthType == reqctx.AuthTypeAdmin {
				spverrors.ErrorResponse(c, spverrors.ErrAdminAuthOnUserEndpoint, reqctx.Logger(c))
				return
			}
		}

		c.Next()
	}
}
