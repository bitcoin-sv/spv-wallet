package data

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/data/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	outpoint, err := bsv.OutpointFromString(c.Param("id"))
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), reqctx.Logger(c))
		return
	}

	data, err := reqctx.Engine(c).DataService().FindForUser(c.Request.Context(), outpoint, userID)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	if data == nil {
		spverrors.ErrorResponse(c, spverrors.ErrDataNotFound, logger)
		return
	}

	c.JSON(http.StatusOK, mapping.DataResponse(data))
}
