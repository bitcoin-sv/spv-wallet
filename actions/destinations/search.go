package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// search will fetch a list of destinations filtered by metadata
// Search Destination godoc
// @Summary		Search for a destination
// @Description	Search for a destination
// @Tags		Destinations
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Destination "List of destinations
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for destinations"
// @Router		/v1/destination/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var destFilters models.DestinationFilters

	destConditions, err := utils.FilterMapByStructFields(*conditions, &destFilters)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var destinations []*engine.Destination
	if destinations, err = a.Services.SpvWalletEngine.GetDestinationsByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		&destConditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := make([]*models.Destination, 0)
	for _, destination := range destinations {
		contracts = append(contracts, mappings.MapToDestinationContract(destination))
	}
	c.JSON(http.StatusOK, contracts)
}
