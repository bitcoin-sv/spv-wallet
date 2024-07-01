package xpubs

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// get will get an existing model
// Get xPub godoc
// @Summary		Get xPub
// @Description	Get xPub
// @Tags		xPub
// @Produce		json
// @Success		200 {object} models.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
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
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindXpub, a.Services.Logger)
		return
	}

	signed := c.GetBool("auth_signed")
	if !signed || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}

// get will get an existing model
// Get xPub godoc
// @Summary		Get xPub
// @Description	Get xPub
// @Tags		xPub
// @Produce		json
// @Success		200 {object} models.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
// @Router		/v1/users/current [get]
// @Security	x-auth-xpub
func (a *Action) get2(c *gin.Context) {
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
