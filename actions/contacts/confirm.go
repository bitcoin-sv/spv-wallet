package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// confirm will confirm contact request
// Confirm contact godoc
// @Summary		Confirm contact - Use (POST) /v1/contacts/{paymail}/confirmation instead
// @Description	This endpoint has been deprecated. Use (POST) /v1/contacts/{paymail}/confirmation instead
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to confirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not unconfirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/confirmed/{paymail} [PATCH]
// @Security	x-auth-xpub
// @Deprecated
func (a *Action) confirm(c *gin.Context) {
	a.confirmContact(c)
}

// confirmContact will confirm contact request
// @Summary		Confirm contact
// @Description	Confirm contact. For contact with status "unconfirmed" change status to "confirmed"
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to confirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not unconfirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contacts/{paymail}/confirmation [POST]
// @Security	x-auth-xpub
func (a *Action) confirmContact(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.ConfirmContact(c, reqXPubID, paymail)

	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}
	c.Status(http.StatusOK)
}
