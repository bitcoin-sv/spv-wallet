package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// oldSearch will fetch a list of utxos filtered on conditions and metadata
// Search UTXO godoc
// @Summary		Search UTXO - Use (GET) /api/v1/utxos instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/utxos instead.
// @Tags		UTXO
// @Produce		json
// @Param		SearchUtxos body filter.SearchUtxos false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Utxo "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchUtxos from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @DeprecatedRouter  /v1/utxo/search [post]
// @Security	x-auth-xpub
func oldSearch(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var reqParams filter.SearchUtxos
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return
	}

	var utxos []*engine.Utxo
	if utxos, err = reqctx.Engine(c).GetUtxosByXpubID(
		c.Request.Context(),
		userContext.GetXPubID(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contracts := make([]*models.Utxo, 0)
	for _, utxo := range utxos {
		contracts = append(contracts, mappings.MapToOldUtxoContract(utxo))
	}

	c.JSON(http.StatusOK, contracts)
}

// search will fetch a list of utxos filtered on conditions and metadata
// Search UTXO godoc
// @Summary		Search UTXO
// @Description	Search UTXO
// @Tags		UTXOs
// @Produce		json
// @Param		SwaggerCommonParams query swagger.CommonFilteringQueryParams false "Supports options for pagination and sorting to streamline data exploration and analysis"
// @Param		UtxoParams query filter.UtxoFilter false "Supports targeted resource searches with filters"
// @Param 		reservedRange[from] query string false "Specifies the start time of the range to query by date when a UTXO was reserved" format(date-time) example:"2024-02-26T11:01:28Z"`
// @Param 		reservedRange[to] query string false "Specifies the end time of the range to query by date when a UTXO was reserved" format(date-time) example:"2024-02-26T11:01:28Z"`
// @Success		200 {object} response.PageModel[response.Utxo] "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchUtxos from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @Router		/api/v1/utxos [get]
// @Security	x-auth-xpub
func search(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)
	searchParams, err := query.ParseSearchParams[filter.UtxoFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams, logger)
		return
	}

	conditions, err := searchParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, logger)
		return
	}

	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	var utxos []*engine.Utxo
	utxos, err = engineInstance.GetUtxosByXpubID(
		c.Request.Context(),
		userContext.GetXPubID(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	utxoContracts := make([]*response.Utxo, 0)
	for _, utxo := range utxos {
		utxoContracts = append(utxoContracts, mappings.MapToUtxoContract(utxo))
	}

	count, err := engineInstance.GetUtxosByXpubIDCount(
		c.Request.Context(),
		userContext.GetXPubID(),
		metadata,
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := response.PageModel[response.Utxo]{
		Content: utxoContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, response)
}
