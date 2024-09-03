package base

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	root := handlersManager.Get(handlers.GroupRoot)
	root.GET("/", index)
	root.OPTIONS("/", statusOK)
	root.HEAD("/", statusOK)

	healthGroup := handlersManager.Group(handlers.GroupRoot, "/"+config.HealthRequestPath)
	healthGroup.GET("", statusOK)
	healthGroup.OPTIONS("", statusOK)
	healthGroup.HEAD("", statusOK)
}

func statusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}
