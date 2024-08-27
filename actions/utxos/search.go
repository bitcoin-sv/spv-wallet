package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
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
func (a *Action) oldSearch(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchUtxos
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, a.Services.Logger)
		return
	}

	var utxos []*engine.Utxo
	if utxos, err = a.Services.SpvWalletEngine.GetUtxosByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
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
// @Tags		UTXO
// @Produce		json
// @Param		SearchUtxos body filter.SearchUtxos false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.Utxo "List of utxos"
// @Failure		400	"Bad request - Error while parsing SearchUtxos from request body"
// @Failure 	500	"Internal server error - Error while searching for utxos"
// @Router		/api/v1/utxos [get]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	searchParams, err := query.ParseSearchParams[filter.SearchUtxos](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams, a.Services.Logger)
		return
	}

	conditions, err := searchParams.Conditions.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, a.Services.Logger)
		return
	}

	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	var utxos []*engine.Utxo
	utxos, err = a.Services.SpvWalletEngine.GetUtxosByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	utxoContracts := make([]*response.Utxo, 0)
	for _, utxo := range utxos {
		utxoContracts = append(utxoContracts, mappings.MapToUtxoContract(utxo))
	}

	count, err := a.Services.SpvWalletEngine.GetUtxosByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := response.PageModel[response.Utxo]{
		Content: utxoContracts,
		Page:    common.GetPageDescriptionFromSearchParams(&searchParams.Page, count),
	}

	c.JSON(http.StatusOK, response)
}
