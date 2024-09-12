package middleware

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// CallbackTokenMiddleware verifies the callback token - if it's valid and matches the Bearer scheme.
func CallbackTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appConfig := reqctx.AppConfig(c)
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrMissingAuthHeader, nil)
		}

		if !strings.HasPrefix(authHeader, BearerSchema) || len(authHeader) <= len(BearerSchema) {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidOrMissingToken, nil)
		}

		providedToken := authHeader[len(BearerSchema):]
		if providedToken != appConfig.Nodes.Callback.Token {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidToken, nil)
		}

		c.Next()
	}
}
