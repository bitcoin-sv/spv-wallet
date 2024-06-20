package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/gin-gonic/gin"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
// Broadcast Callback godoc
// @Summary		Endpoint designed for receiving callbacks from Arc (service responsible for submitting transactions to the BSV network)
// @Tags		Transactions
// @Param 		transaction body broadcast.SubmittedTx true "Transaction"
// @Success		200
// @Failure		400	"Bad request - Error while parsing transaction from request body"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @Router		/transaction/broadcast/callback [post]
// @Security	callback-auth
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
