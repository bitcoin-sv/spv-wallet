package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
// Broadcast Callback godoc
// @Summary		Broadcast Callback
// @Tags		Transactions
// @Param 		transaction body broadcast.SubmittedTx true "transaction"
// @Success		200
// @Router		/transaction/broadcast/callback [post]
// @Security	callback-auth
func (a *Action) broadcastCallback(c *gin.Context) {
	var resp *broadcast.SubmittedTx

	err := c.Bind(&resp)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	err = a.Services.SpvWalletEngine.UpdateTransaction(c.Request.Context(), resp)
	if err != nil {
		a.Services.Logger.Err(err).Msgf("failed to update transaction - tx: %v", resp)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Return response
	c.Status(http.StatusOK)
}
