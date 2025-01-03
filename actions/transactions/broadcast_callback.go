package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
func broadcastCallback(c *gin.Context) {
	logger := reqctx.Logger(c)
	var callbackResp *chainmodels.TXInfo

	err := c.Bind(&callbackResp)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	err = reqctx.Engine(c).HandleTxCallback(c.Request.Context(), callbackResp)
	if err != nil {
		logger.Err(err).Msgf("failed to update transaction - tx: %v", callbackResp)
	}

	c.Status(http.StatusOK)
}
