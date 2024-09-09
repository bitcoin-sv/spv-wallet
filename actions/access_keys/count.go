package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of access keys filtered by metadata
// Count of access keys godoc
// @Summary		Count of access keys - Use (GET) /api/v1/users/current/keys instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current/keys instead.
// @Tags		Access-key
// @Produce		json
// @Param		CountAccessKeys body filter.CountAccessKeys false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountAccessKeys from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of access keys"
// @DeprecatedRouter  /v1/access-key/count [post]
// @Security	x-auth-xpub
func count(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.CountAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetAccessKeysByXPubIDCount(
		c.Request.Context(),
		userContext.GetXPubID(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
