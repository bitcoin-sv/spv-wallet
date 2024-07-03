package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/gin-gonic/gin"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
func (a *Action) broadcastCallback(c *gin.Context) {
	var resp *broadcast.SubmittedTx

	err := c.Bind(&resp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
