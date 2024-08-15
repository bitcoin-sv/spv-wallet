package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// removeContact will confirm contact request
// @Summary		Remove contact
// @Description	Remove contact.
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to confirm"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		500	"Internal server error"
// @Router		/api/v1/contacts/{paymail} [delete]
// @Security	x-auth-xpub
func (a *Action) removeContact(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.DeleteContact(c, reqXPubID, paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.Status(http.StatusOK)
}
