package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
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
func oldRevoke(c *gin.Context, _ *reqctx.UserContext, xpub string) {
	id := c.Query("id")
	revokeHelper(c, id, true, xpub)
}

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key
// @Description	Revoke access key
// @Tags		Access-key
// @Produce		json
// @Param		id path string true "id of the access key"
// @Success		200	{object} response.AccessKey "Revoked AccessKey"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @Router		/api/v1/users/current/keys/{id} [delete]
// @Security	x-auth-xpub
func revoke(c *gin.Context, _ *reqctx.UserContext, xpub string) {
	id := c.Params.ByName("id")
	revokeHelper(c, id, false, xpub)
}

func revokeHelper(c *gin.Context, id string, snakeCase bool, xpub string) {
	logger := reqctx.Logger(c)

	if id == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldID, logger)
		return
	}

	accessKey, err := reqctx.Engine(c).RevokeAccessKey(
		c.Request.Context(),
		xpub,
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
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
