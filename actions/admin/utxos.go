package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// utxosSearch will fetch a list of utxos filtered by metadata
// Search for utxos filtering by metadata godoc
// @Summary		Search for utxos
// @Description	Search for utxos
// @Tags		Admin
// @Produce		json
// @Param		SearchUtxos body filter.AdminSearchUtxos false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Utxo "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchUtxos from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @Router		/v1/admin/utxos/search [post]
// @Security	x-auth-xpub
func utxosSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminSearchUtxos
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return
	}

	utxos, err := reqctx.Engine(c).GetUtxos(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, utxos)
}

// utxosCount will count all utxos filtered by metadata
// Count utxos filtering by metadata godoc
// @Summary		Count utxos
// @Description	Count utxos
// @Tags		Admin
// @Produce		json
// @Param		CountUtxos body filter.AdminCountUtxos false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of utxos"
// @Failure		400	"Bad request - Error while parsing CountUtxos from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of utxos"
// @Router		/v1/admin/utxos/count [post]
// @Security	x-auth-xpub
func utxosCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminCountUtxos
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return
	}

	count, err := reqctx.Engine(c).GetUtxosCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
