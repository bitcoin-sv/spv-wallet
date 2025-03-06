package callback

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

type txSyncService interface {
	Handle(ctx context.Context, txInfo chainmodels.TXInfo) error
}

// RegisterRoutes registers endpoints for callbacks.
func RegisterRoutes(ginEngine *gin.Engine, cfg *config.AppConfig, engine engine.V2Interface) {
	if cfg.ARCCallbackEnabled() {
		callbackURL, err := cfg.ARC.Callback.ShouldGetURL()
		must.HaveNoErrorf(err, "couldn't get callback URL from configuration")

		broadcastCallback := ginEngine.Group("", middleware.CallbackTokenMiddleware())
		broadcastCallback.POST(callbackURL.Path, broadcastCallbackHandler(engine.TxSyncService()))
	}
}

func broadcastCallbackHandler(service txSyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := reqctx.Logger(c)
		var callbackResp chainmodels.TXInfo

		err := c.Bind(&callbackResp)
		if err != nil {
			spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
			return
		}

		err = service.Handle(c, callbackResp)

		if err != nil {
			logger.Err(err).Ctx(c).Any("TxInfo", callbackResp).Msgf("failed to update transaction in ARC broadcast callback handler")
		}

		c.Status(http.StatusOK)
	}
}
