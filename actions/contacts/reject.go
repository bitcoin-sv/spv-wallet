package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// reject will reject contact request
// Reject contact godoc
// @Summary		Reject contact
// @Description	Reject contact. For contact with status "awaiting" delete contact
// @Tags		Contact
// @Produce		json
// @Param		paymail query string true "paymail"
// @Success		200
// @Router		/v1/contact/reject [PATCH]
// @Security	x-auth-xpub
func (a *Action) reject(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.RejectContact(c, reqXPubID, paymail)

	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
