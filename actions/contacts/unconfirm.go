package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

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
func unconfirmContact(c *gin.Context, userContext *reqctx.UserContext) {
	paymail := c.Param("paymail")

	err := reqctx.Engine(c).UnconfirmContact(c.Request.Context(), userContext.GetXPubID(), paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}
	c.Status(http.StatusOK)
}
