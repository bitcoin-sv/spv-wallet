package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// oldUnconfirm will unconfirm contact request
// Unconfirm contact godoc
// @Summary		Unconfirm contact - Use (DELETE) /api/v1/contacts/{paymail}/confirmation instead.
// @Description	This endpoint has been deprecated. Use (DELETE) /api/v1/contacts/{paymail}/confirmation instead.
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to unconfirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not confirmed"
// @Failure		500	"Internal server error"
// @DeprecatedRouter  /v1/contact/unconfirmed/{paymail} [patch]
// @Security	x-auth-xpub
func (a *Action) oldUnconfirm(c *gin.Context) {
	a.unconfirmContact(c)
}

// unconfirmContact will unconfirm contact request
// Unconfirm contact godoc
// @Summary		Unconfirm contact
// @Description	Unconfirm contact. For contact with status "confirmed" change status to "unconfirmed"
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to unconfirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not confirmed"
// @Failure		500	"Internal server error"
// @Router		/api/v1/contacts/{paymail}/confirmation [delete]
// @Security	x-auth-xpub
func (a *Action) unconfirmContact(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.UnconfirmContact(c, reqXPubID, paymail)

	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}
	c.Status(http.StatusOK)
}
