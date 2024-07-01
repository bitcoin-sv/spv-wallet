package contacts

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// confirm will confirm contact request
// Confirm contact godoc
// @Summary		Confirm contact
// @Description	Confirm contact. For contact with status "unconfirmed" change status to "confirmed"
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to confirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not unconfirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/confirmed/{paymail} [PATCH]
// @Security	x-auth-xpub
func (a *Action) confirm(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.ConfirmContact(c, reqXPubID, paymail)

	if err != nil {
		switch {
		case errors.Is(err, engine.ErrContactNotFound):
			c.JSON(http.StatusNotFound, err.Error())
		case errors.Is(err, engine.ErrContactIncorrectStatus):
			c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}

// confirmContact will confirm contact request
// Confirm contact godoc
// @Summary		Confirm contact
// @Description	Confirm contact. For contact with status "unconfirmed" change status to "confirmed"
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to confirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not unconfirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contacts/{paymail}/confirmation [PATCH]
// @Security	x-auth-xpub
func (a *Action) confirmContact(c *gin.Context) {
	a.confirm(c)
}
