package middleware

import (
	"slices"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var securedMiddlewares = []api.MiddlewareFunc{
	(api.MiddlewareFunc)(AuthV2Middleware()),
	(api.MiddlewareFunc)(CheckSignatureMiddleware()),
}

// SignatureAuthWithScopes checks for scopes and runs auth&signature middlewares
func SignatureAuthWithScopes(log *zerolog.Logger) api.MiddlewareFunc {
	return func(c *gin.Context) {
		scopeVal, exists := c.Get(api.XPubAuthScopes)
		if !exists {
			// This means that for this particular endpoint, this auth method is not set
			c.Next()
			return
		}

		scopes, ok := scopeVal.([]string)
		if !ok || len(scopes) == 0 {
			log.Error().Msgf("Invalid scopes for incoming request %s", c.Request.URL.Path)
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInternal, reqctx.Logger(c))
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
				spverrors.AbortWithErrorResponse(c, spverrors.ErrAdminAuthOnUserEndpoint, reqctx.Logger(c))
				return
			}
		case reqctx.AuthTypeXPub:
			if !slices.Contains(scopes, "user") {
				spverrors.AbortWithErrorResponse(c, spverrors.ErrNotAnAdminKey, reqctx.Logger(c))
				return
			}
		default:
			spverrors.AbortWithErrorResponse(c, spverrors.ErrAuthorization, reqctx.Logger(c))
		}

		c.Next()
	}
}
