package accesskeys

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key
// @Description	Revoke access key
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id of the access key"
// @Success		200	{object} models.AccessKey "Revoked AccessKey"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @Router		/v1/access-key [delete]
// @Security	x-auth-xpub
func (a *Action) revoke(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	id := c.Query("id")
	if id == "" {
		spverrors2.ErrorResponse(c, spverrors2.ErrMissingFieldID, a.Services.Logger)
		return
	}

	accessKey, err := a.Services.SpvWalletEngine.RevokeAccessKey(
		c.Request.Context(),
		reqXPub,
		id,
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
