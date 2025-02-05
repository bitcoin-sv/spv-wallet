package middleware

import (
	"slices"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

var securedMiddlewares = []api.MiddlewareFunc{
	(api.MiddlewareFunc)(AuthMiddleware()),
	(api.MiddlewareFunc)(CheckSignatureMiddleware()),
}

// SignatureAuthWithScopes checks for scopes and runs auth&signature middlewares
func SignatureAuthWithScopes() api.MiddlewareFunc {
	return func(c *gin.Context) {
		scopeVal, exists := c.Get(api.XPubAuthScopes)
		if !exists {
			// This means that for this particular endpoint, this auth method is not set
			c.Next()
			return
		}

		scopes, ok := scopeVal.([]string)
		if !ok || len(scopes) == 0 {
			spverrors.ErrorResponse(c, spverrors.ErrWrongAuthScopeFormat, reqctx.Logger(c))
			return
		}

		for _, mid := range securedMiddlewares {
			mid(c)
			if c.IsAborted() {
				return
			}
		}

		userCtx := reqctx.GetUserContext(c)

		switch userCtx.GetAuthType() {
		case reqctx.AuthTypeAdmin:
			if !slices.Contains(scopes, "admin") {
				spverrors.ErrorResponse(c, spverrors.ErrAdminAuthOnNonAdminEndpoint, reqctx.Logger(c))
				return
			}
		case reqctx.AuthTypeXPub:
			if !slices.Contains(scopes, "user") {
				spverrors.ErrorResponse(c, spverrors.ErrUserAuthOnNonUserEndpoint, reqctx.Logger(c))
				return
			}
		default:
			spverrors.ErrorResponse(c, spverrors.ErrAuthorization, reqctx.Logger(c))
		}

		c.Next()
	}
}
