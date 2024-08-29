package middleware

import (
	"net/http"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware is a middleware that handles CORS.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsAllowedHeaders := []string{
			"Content-Type",
			"Cache-Control",
			models.AuthHeader,
			models.AuthAccessKey,
			models.AuthSignature,
			models.AuthHeaderHash,
			models.AuthHeaderNonce,
			models.AuthHeaderTime,
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsAllowedHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
