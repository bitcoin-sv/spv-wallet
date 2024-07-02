package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will count all the utxos fulfilling the given conditions
// Count of UTXOs godoc
// @Summary		Count of UTXOs
// @Description	Count of UTXOs
// @Tags		UTXO
// @Produce		json
// @Param		CountUtxos body filter.CountUtxos false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of utxos"
// @Failure		400	"Bad request - Error while parsing CountUtxos from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of utxos"
// @Router		/v1/utxo/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.CountUtxos
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, a.Services.Logger)
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
		mappings.MapToMetadata(reqParams.Metadata),
		dbConditions,
	); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
