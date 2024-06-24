package destinations

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of destinations filtered by metadata
// Count Destinations godoc
// @Summary		Count Destinations
// @Description	Count Destinations
// @Tags		Destinations
// @Produce		json
// @Param		CountDestinations body filter.CountDestinations false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of destinations"
// @Failure		400	"Bad request - Error while parsing CountDestinations from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of destinations"
// @Router		/v1/destination/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.CountDestinations
	if err := c.Bind(&reqParams); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	count, err := a.Services.SpvWalletEngine.GetDestinationsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
