package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// destinationsSearch will fetch a list of destinations filtered by metadata
// Search for destinations filtering by metadata godoc
// @Summary		Search for destinations
// @Description	Search for destinations
// @Tags		Admin
// @Produce		json
// @Param		SearchDestinations body SearchDestinations false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Destination "List of destinations"
// @Failure		400	"Bad request - Error while parsing SearchDestinations from request body"
// @Failure 	500	"Internal server error - Error while searching for destinations"
// @Router		/v1/admin/destinations/search [post]
// @Security	x-auth-xpub
func (a *Action) destinationsSearch(c *gin.Context) {
	var reqParams SearchDestinations
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	destinations, err := a.Services.SpvWalletEngine.GetDestinations(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	destinationContracts := make([]*models.Destination, 0)
	for _, destination := range destinations {
		destinationContracts = append(destinationContracts, mappings.MapToDestinationContract(destination))
	}

	c.JSON(http.StatusOK, destinationContracts)
}

// destinationsCount will count all destinations filtered by metadata
// Count destinations filtering by metadata godoc
// @Summary		Count destinations
// @Description	Count destinations
// @Tags		Admin
// @Produce		json
// @Param		CountDestinations body CountDestinations false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of destinations"
// @Failure		400	"Bad request - Error while parsing CountDestinations from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of destinations"
// @Security	x-auth-xpub
func (a *Action) destinationsCount(c *gin.Context) {
	var reqParams CountDestinations
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetDestinationsCount(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
