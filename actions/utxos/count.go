package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will count all the utxos fulfilling the given conditions
// Count of UTXOs godoc
// @Summary		Count of UTXOs
// @Description	Count of UTXOs
// @Tags		UTXO
// @Produce		json
// @Param		CountRequestParameters body actions.CountRequestParameters false "Enables precise filtering of resource counts using custom conditions or metadata, catering to specific business or analysis needs"
// @Success		200	{number} int64 "Count of utxos"
// @Failure		400	"Bad request - Error while parsing CountRequestParameters from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of utxos"
// @Router		/v1/utxo/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = conditions
	}

	dbConditions["xpub_id"] = reqXPubID

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetUtxosCount(
		c.Request.Context(),
		metadata,
		dbConditions,
	); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
