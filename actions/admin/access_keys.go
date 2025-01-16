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

// accessKeysSearch will fetch a list of access keys filtered by metadata
// Access Keys Search godoc
// @Summary		Access Keys Search
// @Description	Fetches a list of access keys filtered by metadata, creation range, and other parameters.
// @Tags		Admin
// @Produce		json
// @Param		xpubId query string false "ID of the xPub associated with the access keys"
// @Param		includeDeleted query boolean false "Whether to include deleted access keys"
// @Param		createdRange[from] query string false "Start of creation date range (ISO 8601 format)"
// @Param		createdRange[to] query string false "End of creation date range (ISO 8601 format)"
// @Param		updatedRange[from] query string false "Start of last updated date range (ISO 8601 format)"
// @Param		updatedRange[to] query string false "End of last updated date range (ISO 8601 format)"
// @Param		revokedRange[from] query string false "Start of revoked date range (ISO 8601 format)"
// @Param		revokedRange[to] query string false "End of revoked date range (ISO 8601 format)"
// @Param		page query integer false "Page number for pagination"
// @Param		pageSize query integer false "Number of results per page"
// @Param		orderByField query string false "Field to order results by (e.g., 'created_at')"
// @Param		orderByDirection query string false "Direction of ordering: 'asc' or 'desc'"
// @Success		200 {object} response.PageModel[response.AccessKey] "List of access keys with pagination details"
// @Failure		400 "Bad request - Invalid query parameters"
// @Failure		500 "Internal server error - Error while searching for access keys"
// @Router		/api/v1/admin/users/keys [get]
// @Security	x-auth-xpub
func accessKeysSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	searchParams, err := query.ParseSearchParams[filter.AdminAccessKeyFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions := searchParams.Conditions.ToDbConditions()
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	accessKeys, err := reqctx.Engine(c).GetAccessKeys(
		c.Request.Context(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindAccessKey.WithTrace(err), logger)
		return
	}

	count, err := reqctx.Engine(c).GetAccessKeysCount(c.Request.Context(), metadata, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindAccessKey.WithTrace(err), logger)
		return
	}

	accessKeyContracts := common.MapToTypeContracts(accessKeys, mappings.MapToAccessKeyContract)

	result := response.PageModel[response.AccessKey]{
		Content: accessKeyContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, result)
}
