package accesskeys

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
// @Router		/api/v1/users/current/keys [get]
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
		userContext.GetXPubID(),
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
		userContext.GetXPubID(),
		metadata,
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	res := response.PageModel[response.AccessKey]{
		Content: accessKeyContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, res)
}
