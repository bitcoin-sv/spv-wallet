package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/api"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	root := handlersManager.Get(routes.GroupRoot)

	root.GET("v2/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "v2/swagger/index.html")
	})

	api.Yaml = strings.Replace(api.Yaml, "version: main", fmt.Sprintf("version: '%s'", handlersManager.APIVersion()), 1)
	api.Yaml = strings.Replace(api.Yaml, "https://github.com/bitcoin-sv/spv-wallet/blob/main", fmt.Sprintf("https://github.com/bitcoin-sv/spv-wallet/blob/%s", handlersManager.APIVersion()), 1)

	root.GET("/api/gen.api.yaml", func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml")
		c.String(http.StatusOK, api.Yaml)
	})

	root.GET("v2/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("/api/gen.api.yaml"),
		ginSwagger.PersistAuthorization(true),
		withTitle("SPV Wallet API"),
	))
}

func withTitle(title string) func(*ginSwagger.Config) {
	return func(c *ginSwagger.Config) {
		c.Title = title
	}
}
