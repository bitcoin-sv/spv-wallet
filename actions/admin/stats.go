package admin

import (
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
// @Success		200
// @Router		/v1/admin/stats [get]
// @Security	x-auth-xpub
func (a *Action) stats(c *gin.Context) {
	stats, err := a.Services.SpvWalletEngine.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contract := mappings.MapToAdminStatsContract(stats)
	c.JSON(http.StatusOK, contract)
}
