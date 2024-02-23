package xpubs

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// get will get an existing model
// Get xPub godoc
// @Summary		Get xPub
// @Description	Get xPub
// @Tags		xPub
// @Produce		json
// @Param		key query string false "key"
// @Success		200
// @Router		/v1/xpub [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var xPub *engine.Xpub
	var err error
	if reqXPub != "" {
		xPub, err = a.Services.SpvWalletEngine.GetXpub(
			c.Request.Context(), reqXPub,
		)
	} else {
		xPub, err = a.Services.SpvWalletEngine.GetXpubByID(
			c.Request.Context(), reqXPubID,
		)
	}
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	signed := c.GetBool("auth_signed")
	if !signed || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}
