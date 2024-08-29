package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// oldSearch will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key - Use (GET) /api/v1/users/current/keys instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current/keys instead.
// @Tags		Access-key
// @Produce		json
// @Param		SearchAccessKeys body filter.SearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @DeprecatedRouter  /v1/access-key/search [post]
// @Security	x-auth-xpub
func oldSearch(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var reqParams filter.SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	accessKeys, err := reqctx.Engine(c).GetAccessKeysByXPubID(
		c.Request.Context(),
		userContext.GetXPubID(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	accessKeyContracts := make([]*models.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToOldAccessKeyContract(accessKey))
	}

	c.JSON(http.StatusOK, accessKeyContracts)
}

// search will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key
// @Description	Search access key
// @Tags		Access-key
// @Produce		json
// @Param		SwaggerCommonParams query swagger.CommonFilteringQueryParams false "Supports options for pagination and sorting to streamline data exploration and analysis"
// @Param		AccessKeyParams query filter.AccessKeyFilter false "Supports targeted resource searches with filters"
// @Param 		revokedRange[from] query string false "Specifies the start time of the range to query by date of revoking" format(date-time) example:"2024-02-26T11:01:28Z"`
// @Param 		revokedRange[to] query string false "Specifies the end time of the range to query by date of revoking" format(date-time) example:"2024-02-26T11:01:28Z"`
// @Success		200 {object} response.PageModel[response.AccessKey] "List of access keys"
// @Failure		400	"Bad request - Error while SearchAccessKeys from request query"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @Router		/api/v1//users/current/keys [get]
// @Security	x-auth-xpub
func search(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)

	searchParams, err := query.ParseSearchParams[filter.AccessKeyFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams, logger)
		return
	}

	conditions := searchParams.Conditions.ToDbConditions()
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	accessKeys, err := engine.GetAccessKeysByXPubID(
		c.Request.Context(),
		userContext.GetXPub(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	accessKeyContracts := make([]*response.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	count, err := engine.GetAccessKeysByXPubIDCount(
		c.Request.Context(),
		userContext.GetXPub(),
		metadata,
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := response.PageModel[response.AccessKey]{
		Content: accessKeyContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, response)
}
