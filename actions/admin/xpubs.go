package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		Admin
// @Produce		json
// @Param		CreateXpub body CreateXpub true " "
// @Success		201 {object} models.Xpub "Created Xpub"
// @Failure		400	"Bad request - Error while parsing CreateXpub from request body"
// @Failure 	500	"Internal server error - Error while creating xpub"
// @Router		/v1/admin/xpub [post]
// @Security	x-auth-xpub
func xpubsCreate(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var requestBody CreateXpub
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	xPub, err := reqctx.Engine(c).NewXpub(
		c.Request.Context(), requestBody.Key,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToOldXpubContract(xPub)
	c.JSON(http.StatusCreated, contract)
}

// xpubsSearch will fetch a list of xpubs filtered by metadata
// Search for xpubs filtering by metadata godoc
// @Summary		Search for xpubs
// @Description	Search for xpubs
// @Tags		Admin
// @Produce		json
// @Param		SearchXpubs body filter.SearchXpubs false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Xpub "List of xpubs"
// @Failure		400	"Bad request - Error while parsing SearchXpubs from request body"
// @Failure 	500	"Internal server error - Error while searching for xpubs"
// @Router		/v1/admin/xpubs/search [post]
// @Security	x-auth-xpub
func xpubsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.SearchXpubs
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	xpubs, err := reqctx.Engine(c).GetXPubs(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	xpubContracts := make([]*models.Xpub, 0)
	for _, xpub := range xpubs {
		xpubContracts = append(xpubContracts, mappings.MapToOldXpubContract(xpub))
	}

	c.JSON(http.StatusOK, xpubContracts)
}

// xpubsCount will count all xpubs filtered by metadata
// Count xpubs filtering by metadata godoc
// @Summary		Count xpubs
// @Description	Count xpubs
// @Tags		Admin
// @Produce		json
// @Param		CountXpubs body filter.CountXpubs false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountXpubs from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of xpubs"
// @Router		/v1/admin/xpubs/count [post]
// @Security	x-auth-xpub
func xpubsCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.CountXpubs
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetXPubsCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
