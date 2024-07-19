package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// unconfirm will unconfirm contact request
// Unconfirm contact godoc
// @Summary		Unconfirm contact
// @Description	Unconfirm contact. For contact with status "confirmed" change status to "unconfirmed"
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to unconfirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not confirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/unconfirmed/{paymail} [PATCH]
// @Security	x-auth-xpub
func (a *Action) unconfirm(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.UnconfirmContact(c, reqXPubID, paymail)

	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}
	c.Status(http.StatusOK)
}

// unconfirmContact will unconfirm contact request
// Unconfirm contact godoc
// @Summary		Unconfirm contact
// @Description	Unconfirm contact. For contact with status "confirmed" change status to "unconfirmed"
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to unconfirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not confirmed"
// @Failure		500	"Internal server error"
// @Router		/v1/contacts/{paymail}/non-confirmation [PATCH]
// @Security	x-auth-xpub
func (a *Action) unconfirmContact(c *gin.Context) {
	a.unconfirm(c)
}
