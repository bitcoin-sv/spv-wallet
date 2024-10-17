package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// destinationsSearch will fetch a list of destinations filtered by metadata
// Search for destinations filtering by metadata godoc
// @DDeprecatedRouter /v1/admin/destinations [post]
// @Summary		Search for destinations
// @Description	Search for destinations
// @Tags		Admin
// @Produce		json
// @Param		SearchDestinations body filter.SearchDestinations false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Destination "List of destinations"
// @Failure		400	"Bad request - Error while parsing SearchDestinations from request body"
// @Failure 	500	"Internal server error - Error while searching for destinations"
// @Router		/v1/admin/destinations/search [post]
// @Security	x-auth-xpub
func destinationsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.SearchDestinations
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	destinations, err := reqctx.Engine(c).GetDestinations(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	destinationContracts := make([]*models.Destination, 0)
	for _, destination := range destinations {
		destinationContracts = append(destinationContracts, mappings.MapOldToDestinationContract(destination))
	}

	c.JSON(http.StatusOK, destinationContracts)
}

// destinationsCount will count all destinations filtered by metadata
// Count destinations filtering by metadata godoc
// @DeprecatedRouter /v1/admin/destinations/count [post]
// @Summary		Count destinations
// @Description	Count destinations
// @Tags		Admin
// @Produce		json
// @Param		CountDestinations body filter.CountDestinations false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of destinations"
// @Failure		400	"Bad request - Error while parsing CountDestinations from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of destinations"
// @Security	x-auth-xpub
func destinationsCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.CountDestinations
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetDestinationsCount(
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
