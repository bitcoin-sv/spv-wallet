package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
func (a *Action) broadcastCallback(c *gin.Context) {
	var resp *broadcast.SubmittedTx

	err := c.Bind(&resp)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	err = a.Services.SpvWalletEngine.UpdateTransaction(c.Request.Context(), resp)
	if err != nil {
		a.Services.Logger.Err(err).Msgf("failed to update transaction - tx: %v", resp)
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	// Return response
	c.Status(http.StatusOK)
}
