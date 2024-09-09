package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of transactions filtered on conditions and metadata
// Count of transactions godoc
// @Summary		Count of transactions - Use (GET) /api/v1/transactions instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/transactions instead
// @Tags		Transactions
// @Produce		json
// @Param		CountTransactions body filter.CountTransactions false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountTransactions from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of transactions"
// @DeprecatedRouter		/v1/transaction/count [post]
// @Security	x-auth-xpub
func count(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var reqParams filter.CountTransactions
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetTransactionsByXpubIDCount(
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
