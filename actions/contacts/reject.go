package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// oldReject will reject contact request
// Reject contact godoc
// @Summary		Reject contact - Use (DELETE) /api/v1/invitations/{paymail} instead.
// @Description	This endpoint has been deprecated. Use (DELETE) /api/v1/invitations/{paymail} instead.
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to reject"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not awaiting"
// @Failure		500	"Internal server error"
// @DeprecatedRouter  /v1/contact/rejected/{paymail} [patch]
// @Security	x-auth-xpub
func oldReject(c *gin.Context, userContext *reqctx.UserContext) {
	rejectInvitation(c, userContext)
}

// rejectInvitation will reject contact request
// Reject contact invitation godoc
// @Summary		Reject contact invitation
// @Description	Reject contact invitation. For contact with status "awaiting" delete contact
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to reject"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		422	"Contact status not awaiting"
// @Failure		500	"Internal server error"
// @Router		/api/v1/invitations/{paymail} [delete]
// @Security	x-auth-xpub
func rejectInvitation(c *gin.Context, userContext *reqctx.UserContext) {
	paymail := c.Param("paymail")

	err := reqctx.Engine(c).RejectContact(c, userContext.GetXPubID(), paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}
	c.Status(http.StatusOK)
}
