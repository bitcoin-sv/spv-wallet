package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key - Use (DELETE) /api/v1/users/current/keys/{id} instead.
// @Description	This endpoint has been deprecated. Use (DELETE) /api/v1/users/current/keys/{id} instead.
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id of the access key"
// @Success		200	{object} models.AccessKey "Revoked AccessKey"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @DeprecatedRouter  /v1/access-key [delete]
// @Security	x-auth-xpub
func (a *Action) oldRevoke(c *gin.Context) {
	a.revokeHelper(c, true)
}

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key
// @Description	Revoke access key
// @Tags		Access-key
// @Produce		json
// @Param		id path string true "id of the access key"
// @Success		200	{object} models.AccessKey "Revoked AccessKey"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @Router		/api/v1/users/current/keys/{id} [delete]
// @Security	x-auth-xpub
func (a *Action) revoke(c *gin.Context) {
	a.revokeHelper(c, false)
}

func (a *Action) revokeHelper(c *gin.Context, snakeCase bool) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	id := c.Params.ByName("id")
	if id == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldID, a.Services.Logger)
		return
	}

	accessKey, err := a.Services.SpvWalletEngine.RevokeAccessKey(
		c.Request.Context(),
		reqXPub,
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	if snakeCase {
		contract := mappings.MapToOldAccessKeyContract(accessKey)
		c.JSON(http.StatusCreated, contract)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
