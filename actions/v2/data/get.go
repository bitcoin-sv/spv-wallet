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

	dataID := c.Param("id")
	_, err = bsv.OutpointFromString(dataID)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidDataID.Wrap(err), logger)
		return
	}

	data, err := reqctx.Engine(c).DataService().FindForUser(c.Request.Context(), dataID, userID)
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
