package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
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
func removeContact(c *gin.Context, userContext *reqctx.UserContext) {
	paymail := c.Param("paymail")

	err := reqctx.Engine(c).DeleteContact(c.Request.Context(), userContext.GetXPubID(), paymail)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.Status(http.StatusOK)
}
