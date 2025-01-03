package paymailserver

import (
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/gin-gonic/gin"
)

// Register registers the paymail server.
func Register(configuration *server.Configuration, ginEngine *gin.Engine) {
	configuration.RegisterRoutes(ginEngine)
}
