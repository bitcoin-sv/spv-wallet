package admin

import (
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/gin-gonic/gin"
)

// status will return the status of the admin login
// Get stats godoc
// @Summary		Get stats
// @Description	Get stats
// @Tags		Admin
// @Produce		json
// @Success		200	{object} models.AdminStats "Stats for the admin"
// @Failure 	500	"Internal Server Error - Error while fetching admin stats"
// @Router		/v1/admin/stats [get]
// @Security	x-auth-xpub
func (a *Action) stats(c *gin.Context) {
	stats, err := a.Services.SpvWalletEngine.GetStats(c.Request.Context())
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToAdminStatsContract(stats)
	c.JSON(http.StatusOK, contract)
}
