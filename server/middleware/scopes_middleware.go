package middleware

import (
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
			spverrors.ErrorResponse(c, spverrors.ErrMissingAuthScope, reqctx.Logger(c))
			return
		}

		scopes, ok := scopeVal.([]string)
		if !ok || len(scopes) == 0 {
			spverrors.ErrorResponse(c, spverrors.ErrWrongAuthScopeFormat, reqctx.Logger(c))
			return
		}

		userType := getHighestPriorityScope(scopes)

		if userType == "" {
			spverrors.ErrorResponse(c, spverrors.ErrWrongAuthScopeType, reqctx.Logger(c))
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

		c.Next()
	}
}

// getHighestPriorityScope returns the highest priority scope if many were defined
func getHighestPriorityScope(scopes []string) string {
	scopePriority := map[string]int{
		"admin": 3,
		"user":  2,
		"basic": 1,
	}

	highestScope := ""
	highestPriority := 0

	for _, scope := range scopes {
		if priority, exists := scopePriority[scope]; exists && priority > highestPriority {
			highestScope = scope
			highestPriority = priority
		}
	}

	return highestScope
}
