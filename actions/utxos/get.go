package utxos

import (
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
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
func (a *Action) get(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	txID := c.Query("tx_id")
	outputIndex := c.Query("output_index")
	outputIndex64, err := strconv.ParseUint(outputIndex, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	utxo, err := a.Services.SpvWalletEngine.GetUtxo(
		c.Request.Context(),
		reqXPubID,
		txID,
		uint32(outputIndex64),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToUtxoContract(utxo)
	c.JSON(http.StatusOK, contract)
}
