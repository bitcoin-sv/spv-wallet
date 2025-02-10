package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// Register collects all the action's routes and registers them using the handlersManager
func Register(handlersManager *handlers.Manager) {
	transactions.RegisterRoutes(handlersManager)

	admin.Register(handlersManager)

	RegisterRoutes(handlersManager)
}

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	root := handlersManager.Get(handlers.GroupRoot)

	root.GET("v2/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "v2/swagger/index.html")
	})
	root.StaticFile("/api/gen.api.yaml", "./api/gen.api.yaml")
	root.GET("v2/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/api/gen.api.yaml")))
}
