package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// @Summary				Get stats Use (GET) /api/v1/admin/stats instead.
// @Description			This endpoint has been deprecated. Use (GET) /api/v1/admin/stats instead.
// @Tags				Admin
// @Produce				json
// @Success				200	{object} models.AdminStats "Stats for the admin"
// @Failure 			500	"Internal Server Error - Error while fetching admin stats"
// @DeprecatedRouter	/v1/admin/stats [get]
// @Security			x-auth-xpub
func statsOld(c *gin.Context, _ *reqctx.AdminContext) {
	stats, err := reqctx.Engine(c).GetStats(c.Request.Context())
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	contract := mappings.MapToOldAdminStatsContract(stats)
	c.JSON(http.StatusOK, contract)
}

// @Summary			Get stats
// @Description		Get statistics of the spv-wallet
// @Tags			Admin
// @Produce			json
// @Success			200	{object} response.AdminStats "Stats for the admin"
// @Failure 		500	"Internal Server Error - Error while fetching admin stats"
// @Router			/api/v1/admin/stats [get]
// @Security		x-auth-xpub
func stats(c *gin.Context, _ *reqctx.AdminContext) {
	stats, err := reqctx.Engine(c).GetStats(c.Request.Context())
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	contract := mappings.MapToAdminStatsContract(stats)
	c.JSON(http.StatusOK, contract)
}
