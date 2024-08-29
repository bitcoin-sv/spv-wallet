package utxos

import (
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// get will fetch a given utxo according to conditions
// Get UTXO godoc
// @Summary		Get UTXO
// @Description	Get UTXO
// @Tags		UTXO
// @Produce		json
// @Param		tx_id query string true "Id of the transaction"
// @Param		output_index query int true "Output index"
// @Success		200 {object} models.Utxo "UTXO with given Id and output index"
// @Failure		400	"Bad request - Error while parsing output_index"
// @Failure 	500	"Internal Server Error - Error while fetching utxo"
// @Router		/v1/utxo [get]
// @Security	x-auth-xpub
func get(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	txID := c.Query("tx_id")
	outputIndex := c.Query("output_index")
	outputIndex64, err := strconv.ParseUint(outputIndex, 10, 32)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	utxo, err := reqctx.Engine(c).GetUtxo(
		c.Request.Context(),
		userContext.GetXPubID(),
		txID,
		uint32(outputIndex64),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToUtxoContract(utxo)
	c.JSON(http.StatusOK, contract)
}
