package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		Admin
// @Produce		json
// @Param		CreateXpub body CreateXpub true " "
// @Success		201 {object} response.Xpub "Created Xpub"
// @Failure		400	"Bad request - Error while parsing CreateXpub from request body"
// @Failure 	500	"Internal server error - Error while creating xpub"
// @Router		/api/v1/admin/users [post]
// @Security	x-auth-xpub
func xpubsCreate(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var requestBody CreateXpub
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	xPub, err := reqctx.Engine(c).NewXpub(
		c.Request.Context(), requestBody.Key,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidRequesterXpub.WithTrace(err), logger)
		return
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusCreated, contract)
}

// xpubsSearch will fetch a list of xpubs filtered by metadata
// Search for xpubs filtering by metadata godoc
// @Summary		Search for xpubs
// @Description	Fetches a list of xpubs filtered by metadata and other criteria
// @Tags		Admin
// @Produce		json
// @Param		SwaggerCommonParams query swagger.CommonFilteringQueryParams false "Supports options for pagination and sorting to streamline data exploration and analysis"
// @Param		XpubFilter query filter.XpubFilter false "Supports targeted resource searches with filters"
// @Param		id query string false "XPub ID (UUID)"
// @Param		currentBalance query integer false "Current balance of the xPub"
// @Success		200 {object} response.PageModel[response.Xpub] "List of xPubs with pagination details"
// @Failure		400 "Bad request - Invalid query parameters"
// @Failure		500 "Internal server error - Error while searching for xPubs"
// @Router 		/api/v1/admin/users [get]
// @Security	x-auth-xpub
func xpubsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.XpubFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	xpubs, err := reqctx.Engine(c).GetXPubs(
		c.Request.Context(),
		mappings.MapToMetadata(searchParams.Metadata),
		searchParams.Conditions.ToDbConditions(),
		mappings.MapToDbQueryParams(&searchParams.Page),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	count, err := reqctx.Engine(c).GetXPubsCount(c.Request.Context(), mappings.MapToMetadata(searchParams.Metadata), searchParams.Conditions.ToDbConditions())
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotCountXpubs.WithTrace(err), logger)
		return
	}

	xpubContracts := common.MapToTypeContracts(xpubs, mappings.MapToXpubContract)

	result := response.PageModel[response.Xpub]{
		Content: xpubContracts,
		Page:    common.GetPageDescriptionFromSearchParams(mappings.MapToDbQueryParams(&searchParams.Page), count),
	}

	c.JSON(http.StatusOK, result)
}
