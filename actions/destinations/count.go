package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of destinations filtered by metadata
// Count Destinations godoc
// @Summary		Count Destinations. This endpoint has been deprecated (it will be removed in the future).
// @Description	Count Destinations. This endpoint has been deprecated (it will be removed in the future).
// @Tags		Destinations
// @Produce		json
// @Param		CountDestinations body filter.CountDestinations false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of destinations"
// @Failure		400	"Bad request - Error while parsing CountDestinations from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of destinations"
// @DeprecatedRouter  /v1/destination/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.CountDestinations
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	count, err := a.Services.SpvWalletEngine.GetDestinationsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
