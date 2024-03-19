package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// accept will accept contact request
// Accept contact godoc
// @Summary		Accept contact
// @Description	Accept contact. For contact with status "awaiting" change status to "unconfirmed"
// @Tags		Contact
// @Produce		json
// @Param		paymail query string true "paymail"
// @Success		200
// @Router		/v1/contact/accept [PATCH]
// @Security	x-auth-xpub
func (a *Action) accept(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	err := a.Services.SpvWalletEngine.AcceptContact(c, reqXPubID, paymail)

	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, "contact accepted")
}
