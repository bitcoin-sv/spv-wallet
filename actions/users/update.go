package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// oldUpdate will update an existing model
// Update current user information godoc
// @Summary		Update current user information - Use (PATCH) /api/v1/users/current instead.
// @Description	This endpoint has been deprecated. Use (PATCH) /api/v1/users/current instead.
// @Tags		Users
// @Produce		json
// @Param		Metadata body engine.Metadata false " "
// @Success		200 {object} models.Xpub "Updated xPub"
// @Failure		400	"Bad request - Error while parsing Metadata from request body"
// @Failure 	500	"Internal Server Error - Error while updating xPub"
// @DeprecatedRouter  /v1/xpub [patch]
// @Security	x-auth-xpub
func (a *Action) oldUpdate(c *gin.Context) {
	a.updateHelper(c, true)
}

// update will update an existing model
// Update current user information godoc
// @Summary		Update current user information
// @Description	Update current user information
// @Tags		Users
// @Produce		json
// @Param		Metadata body engine.Metadata false " "
// @Success		200 {object} response.Xpub "Updated xPub"
// @Failure		400	"Bad request - Error while parsing Metadata from request body"
// @Failure 	500	"Internal Server Error - Error while updating xPub"
// @Router		/api/v1/users/current [patch]
// @Security	x-auth-xpub
func (a *Action) update(c *gin.Context) {
	a.updateHelper(c, false)
}

func (a *Action) updateHelper(c *gin.Context, snakeCase bool) {
	reqXPub := c.GetString(auth.ParamXPubKey)
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var requestBody engine.Metadata
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	// Get an xPub
	var xPub *engine.Xpub
	var err error
	xPub, err = a.Services.SpvWalletEngine.UpdateXpubMetadata(
		c.Request.Context(), reqXPubID, requestBody,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	signed := c.GetBool("auth_signed")
	if !signed || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	if snakeCase {
		contract := mappings.MapToOldXpubContract(xPub)
		c.JSON(http.StatusOK, contract)
		return
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}
