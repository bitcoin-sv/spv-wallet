package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of destinations filtered by metadata
// Count Destinations godoc
// @Summary		Count Destinations
// @Description	Count Destinations
// @Tags		Destinations
// @Produce		json
// @Param		CountDestinations body CountDestinations false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of destinations"
// @Failure		400	"Bad request - Error while parsing CountDestinations from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of destinations"
// @Router		/v1/destination/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams CountDestinations
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetDestinationsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
