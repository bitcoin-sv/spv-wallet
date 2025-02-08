package base

import (
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	root := handlersManager.Get(routes.GroupRoot)

	root.GET("v2/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "v2/swagger/index.html")
	})
	root.StaticFile("/api/gen.api.yaml", "./api/gen.api.yaml")
	root.GET("v2/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/api/gen.api.yaml")))
}
