package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
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
// @Success		200	{object} models.AccessKey "Created AccessKey"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @Router		/v1/access-key [delete]
// @Security	x-auth-xpub
func (a *Action) revoke(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, engine.ErrMissingFieldID)
		return
	}

	accessKey, err := a.Services.SpvWalletEngine.RevokeAccessKey(
		c.Request.Context(),
		reqXPub,
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
