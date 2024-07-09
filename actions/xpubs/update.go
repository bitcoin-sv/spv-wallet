package xpubs

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// update will update an existing model
// Update xPub godoc
// @Summary		Update xPub
// @Description	Update xPub
// @Tags		xPub
// @Produce		json
// @Param		Metadata body engine.Metadata false " "
// @Success		200 {object} models.Xpub "Updated xPub"
// @Failure		400	"Bad request - Error while parsing Metadata from request body"
// @Failure 	500	"Internal Server Error - Error while updating xPub"
// @Router		/v1/xpub [patch]
// @Security	x-auth-xpub
// @Security	bux-auth-xpub
func (a *Action) update(c *gin.Context) {
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

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusOK, contract)
}
