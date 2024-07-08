package xpubs

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// oldGet will get an existing model
// Get current user information godoc
// @Summary		Get current user information - Use (GET) /api/v1/users/current instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current instead.
// @Tags		Users
// @Produce		json
// @Success		200 {object} models.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
// @DeprecatedRouter  /v1/xpub [get]
// @Security	x-auth-xpub
func (a *Action) oldGet(c *gin.Context) {
	a.get(c)
}

// get will get an existing model
// Get current user information godoc
// @Summary		Get current user information
// @Description	Get current user information
// @Tags		Users
// @Produce		json
// @Success		200 {object} models.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
// @Router		/api/v1/users/current [get]
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	signed := c.GetBool("auth_signed")
	if !signed || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}
