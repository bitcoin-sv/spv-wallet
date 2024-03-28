package contacts

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// reject will reject contact request
// Reject contact godoc
// @Summary		Reject contact
// @Description	Reject contact. For contact with status "awaiting" delete contact
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to reject"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not awaiting"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/rejected [PATCH]
// @Security	x-auth-xpub
func (a *Action) reject(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.RejectContact(c, reqXPubID, paymail)

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
