package contacts

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// accept will accept contact request
// Accept contact godoc
// @Summary		Accept contact
// @Description	Accept contact. For contact with status "awaiting" change status to "unconfirmed"
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to accept"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not awaiting"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/accepted [PATCH]
// @Security	x-auth-xpub
func (a *Action) accept(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.AcceptContact(c, reqXPubID, paymail)

	if err != nil {
		switch {
		case errors.Is(err, engine.ErrContactNotFound):
			c.JSON(http.StatusNotFound, err.Error())
		case errors.Is(err, engine.ErrContactStatusNotAwaiting):
			c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}