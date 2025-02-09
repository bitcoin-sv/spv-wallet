package data

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/data/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIData) GetApiV2DataId(c *gin.Context, id string) {
	logger := reqctx.Logger(c)

	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	_, err = bsv.OutpointFromString(id)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidDataID.Wrap(err), logger)
		return
	}

	data, err := reqctx.Engine(c).DataService().FindForUser(c.Request.Context(), id, userID)
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
