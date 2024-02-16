package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
)

// count will fetch a count of destinations filtered by metadata
// Count Destinations godoc
// @Summary		Count Destinations
// @Description	Count Destinations
// @Tags		Destinations
// @Param		metadata query string false "metadata"
// @Param		condition query string false "condition"
// @Produce		json
// @Success		200
// @Router		/v1/destination/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	_, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)
	var count int64
	if count, err = a.Services.SpvWalletEngine.GetDestinationsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
