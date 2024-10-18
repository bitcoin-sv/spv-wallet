package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
)

// utxosSearch will fetch a list of utxos filtered by metadata
// Search for utxos filtering by metadata godoc
// @Summary		Search for utxos
// @Description	Search for utxos
// @Tags		Admin
// @Produce		json
// @Param		SearchUtxos body filter.AdminUtxoFilter false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.Utxo "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchUtxos from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @Router		/api/v1/admin/utxos [get]
// @Security	x-auth-xpub
func utxosSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.AdminUtxoFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions, err := searchParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions.WithTrace(err), logger)
		return
	}
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	utxos, err := reqctx.Engine(c).GetUtxos(
		c.Request.Context(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindUtxo.WithTrace(err), logger)
		return
	}

	contracts := make([]*response.Utxo, 0, len(utxos))
	for _, utxo := range utxos {
		contracts = append(contracts, mappings.MapToUtxoContract(utxo))
	}

	c.JSON(http.StatusOK, contracts)
}
