package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
// @Router		/v1/contact/rejected/{paymail} [PATCH]
// @Security	x-auth-xpub
func (a *Action) reject(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.RejectContact(c, reqXPubID, paymail)

	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}
	c.Status(http.StatusOK)
}

// rejectInvitation will reject contact request
// Reject contact invitation godoc
// @Summary		Reject contact invitation
// @Description	Reject contact invitation. For contact with status "awaiting" delete contact
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to reject"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not awaiting"
// @Failure		500	"Internal server error"
// @Router		/v1/invitations/{paymail} [DELETE]
// @Security	x-auth-xpub
func (a *Action) rejectInvitation(c *gin.Context) {
	a.reject(c)
}
