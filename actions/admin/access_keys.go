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

// accessKeysSearch will fetch a list of access keys filtered by metadata
// Access Keys Search godoc
// @Summary		Access Keys Search
// @Description	Access Keys Search
// @Tags		Admin
// @Produce		json
// @Param		SearchAccessKeys body filter.AdminSearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while parsing SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
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

	accessKeyContracts := make([]*response.AccessKey, 0, len(accessKeys))
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	c.JSON(http.StatusOK, accessKeyContracts)
}
