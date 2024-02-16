package xpubs

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// update will update an existing model
// Update xPub godoc
// @Summary		Update xPub
// @Description	Update xPub
// @Tags		xPub
// @Produce		json
// @Param		metadata query string false "metadata"
// @Success		200
// @Router		/v1/xpub [patch]
// @Security	x-auth-xpub
// @Security	bux-auth-xpub
func (a *Action) update(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var requestBody engine.Metadata
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get an xPub
	var xPub *engine.Xpub
	var err error
	xPub, err = a.Services.SpvWalletEngine.UpdateXpubMetadata(
		c.Request.Context(), reqXPubID, requestBody,
	)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
	}

	signed := c.GetBool("auth_signed")
	if signed == false || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}
