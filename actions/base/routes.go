package base

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/docs"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	docs.SwaggerInfo.Version = handlersManager.APIVersion()
	root := handlersManager.Get(handlers.GroupRoot)
	root.GET("/", index)
	root.OPTIONS("/", statusOK)
	root.HEAD("/", statusOK)

	healthGroup := handlersManager.Group(handlers.GroupRoot, "/"+config.HealthRequestPath)
	healthGroup.GET("", statusOK)
	healthGroup.OPTIONS("", statusOK)
	healthGroup.HEAD("", statusOK)

	root.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func statusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}
