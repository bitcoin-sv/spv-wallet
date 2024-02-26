package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
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
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "metadata"
// @Param		condition query string false "condition"
// @Success		200
// @Router		/v1/destination/search [get]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var destinations []*engine.Destination
	if destinations, err = a.Services.SpvWalletEngine.GetDestinationsByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contracts := make([]*models.Destination, 0)
	for _, destination := range destinations {
		contracts = append(contracts, mappings.MapToDestinationContract(destination))
	}
	c.JSON(http.StatusOK, contracts)
}
