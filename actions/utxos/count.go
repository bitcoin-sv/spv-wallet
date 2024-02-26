package utxos

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// count will count all the utxos fulfilling the given conditions
// Count of UTXOs godoc
// @Summary		Count of UTXOs
// @Description	Count of UTXOs
// @Tags		UTXO
// @Produce		json
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/utxo/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	_, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = *conditions
	}

	dbConditions["xpub_id"] = reqXPubID

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetUtxosCount(
		c.Request.Context(),
		metadata,
		&dbConditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
