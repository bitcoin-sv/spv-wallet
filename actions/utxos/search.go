package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// search will fetch a list of utxos filtered on conditions and metadata
// Search UTXO godoc
// @Summary		Search UTXO
// @Description	Search UTXO
// @Tags		UTXO
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Utxo "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @Router		/v1/utxo/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var utxos []*engine.Utxo
	if utxos, err = a.Services.SpvWalletEngine.GetUtxosByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := make([]*models.Utxo, 0)
	for _, utxo := range utxos {
		contracts = append(contracts, mappings.MapToUtxoContract(utxo))
	}

	c.JSON(http.StatusOK, contracts)
}
