package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// utxosSearch will fetch a list of utxos filtered by metadata
// Search for utxos filtering by metadata godoc
// @Summary		Search for utxos
// @Description	Fetches a list of UTXOs filtered by metadata and other criteria
// @Tags		Admin
// @Produce		json
// @Param		id query string false "UTXO ID (UUID)"
// @Param		transactionId query string false "Transaction ID associated with the UTXO"
// @Param		outputIndex query integer false "Output index of the UTXO"
// @Param		satoshis query integer false "Amount of satoshis held in the UTXO"
// @Param		scriptPubKey query string false "ScriptPubKey associated with the UTXO"
// @Param		type query string false "Type of the UTXO (e.g., 'P2PKH', 'P2SH')"
// @Param		draftId query string false "Draft ID associated with the UTXO"
// @Param		reservedRange[from] query string false "Start of reserved date range (ISO 8601 format)"
// @Param		reservedRange[to] query string false "End of reserved date range (ISO 8601 format)"
// @Param		spendingTxId query string false "Transaction ID spending the UTXO"
// @Param		xpubId query string false "XPub ID associated with the UTXO"
// @Param		includeDeleted query boolean false "Whether to include deleted UTXOs"
// @Param		createdRange[from] query string false "Start of creation date range (ISO 8601 format)"
// @Param		createdRange[to] query string false "End of creation date range (ISO 8601 format)"
// @Param		updatedRange[from] query string false "Start of last updated date range (ISO 8601 format)"
// @Param		updatedRange[to] query string false "End of last updated date range (ISO 8601 format)"
// @Param		page query integer false "Page number for pagination"
// @Param		pageSize query integer false "Number of results per page"
// @Param		orderByField query string false "Field to order results by (e.g., 'created_at')"
// @Param		orderByDirection query string false "Direction of ordering: 'asc' or 'desc'"
// @Success		200 {object} response.PageModel[response.Utxo] "List of UTXOs with pagination details"
// @Failure		400 "Bad request - Invalid query parameters"
// @Failure		500 "Internal server error - Error while searching for UTXOs"
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

	count, err := reqctx.Engine(c).GetUtxosCount(c.Request.Context(), metadata, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindUtxo.WithTrace(err), logger)
		return
	}

	contracts := common.MapToTypeContracts(utxos, mappings.MapToUtxoContract)

	result := response.PageModel[response.Utxo]{
		Content: contracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, result)
}
