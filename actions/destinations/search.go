package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// search will fetch a list of destinations filtered by metadata
// Search Destination godoc
// @Summary		Search for a destination
// @Description	Search for a destination
// @Tags		Destinations
// @Produce		json
// @Param		SearchDestinations body SearchDestinations false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Destination "List of destinations
// @Failure		400	"Bad request - Error while parsing SearchDestinations from request body"
// @Failure 	500	"Internal server error - Error while searching for destinations"
// @Router		/v1/destination/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchDestinations
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	destinations, err := a.Services.SpvWalletEngine.GetDestinationsByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := make([]*models.Destination, 0)
	for _, destination := range destinations {
		contracts = append(contracts, mappings.MapToDestinationContract(destination))
	}
	c.JSON(http.StatusOK, contracts)
}
