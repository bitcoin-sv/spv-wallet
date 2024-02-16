package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// destinationsSearch will fetch a list of destinations filtered by metadata
// Search for destinations filtering by metadata godoc
// @Summary		Search for destinations
// @Description	Search for destinations
// @Tags		Admin
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/destinations/search [post]
// @Security	x-auth-xpub
func (a *Action) destinationsSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var destinations []*engine.Destination
	if destinations, err = a.Services.SpvWalletEngine.GetDestinations(
		c.Request.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, destinations)
}

// destinationsCount will count all destinations filtered by metadata
// Count destinations filtering by metadata godoc
// @Summary		Count destinations
// @Description	Count destinations
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/destinations/count [post]
// @Security	x-auth-xpub
func (a *Action) destinationsCount(c *gin.Context) {
	_, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetDestinationsCount(
		c.Request.Context(),
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
