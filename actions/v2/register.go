package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/paymailserver"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/callback"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/swagger"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
)

// RegisterNonOpenAPIRoutes collects all the action's routes that aren't part of the Open API documentation and registers them using the handlersManager.
func RegisterNonOpenAPIRoutes(ginEngine *gin.Engine, cfg *config.AppConfig, engine engine.V2Interface) {
	paymailserver.Register(engine.PaymailServerConfiguration(), ginEngine)
	swagger.RegisterRoutes(ginEngine, cfg)
	callback.RegisterRoutes(ginEngine, cfg, engine)
}
