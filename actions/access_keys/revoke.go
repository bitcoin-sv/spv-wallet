package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key
// @Description	Revoke access key
// @Tags		Access-key
// @Produce		json
// @Param		id path string true "id of the access key"
// @Success		200
// @Failure		400	"Bad request - Missing required field: id"
// @Failure 	500	"Internal server error - Error while revoking access key"
// @Router		/api/v1/users/current/keys/{id} [delete]
// @Security	x-auth-xpub
func revoke(c *gin.Context, userContext *reqctx.UserContext) {
	id := c.Params.ByName("id")
	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		return
	}
	logger := reqctx.Logger(c)

	if id == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldID, logger)
		return
	}

	_, err = reqctx.Engine(c).RevokeAccessKey(
		c.Request.Context(),
		xpub,
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.Status(http.StatusOK)
}
