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
	root := handlersManager.Get(handlers.GroupRoot)
	root.GET("/", index)
	root.OPTIONS("/", statusOK)
	root.HEAD("/", statusOK)

	healthGroup := handlersManager.Group(handlers.GroupRoot, "/"+config.HealthRequestPath)
	healthGroup.GET("", statusOK)
	healthGroup.OPTIONS("", statusOK)
	healthGroup.HEAD("", statusOK)

	// Register Swagger for v1 API
	docs.SwaggerInfo.Version = handlersManager.APIVersion()
	root.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Register Swagger for v2 API
	root.GET("v2/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "v2/swagger/index.html")
	})
	root.StaticFile("/docs/openapi_v2.yaml", "./docs/openapi_v2.yaml")
	root.GET("v2/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/docs/openapi_v2.yaml")))

}

func statusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}
