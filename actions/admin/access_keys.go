package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// accessKeysSearch will fetch a list of access keys filtered by metadata
// Access Keys Search godoc
// @Summary		Access Keys Search
// @Description	Access Keys Search
// @Tags		Admin
// @Produce		json
// @Param		SearchAccessKeys body filter.AdminSearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while parsing SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @Router		/v1/admin/access-keys/search [post]
// @Security	x-auth-xpub
func accessKeysSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminSearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	accessKeys, err := reqctx.Engine(c).GetAccessKeys(
		c.Request.Context(),
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

// accessKeysCount will count all access keys filtered by metadata
// Access Keys Count godoc
// @Summary		Access Keys Count
// @Description	Access Keys Count
// @Tags		Admin
// @Produce		json
// @Param		CountAccessKeys body filter.AdminCountAccessKeys false "Enables filtering of elements to be counted"
// @Success		200 {number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountAccessKeys from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of access keys"
// @Router		/v1/admin/access-keys/count [post]
// @Security	x-auth-xpub
func accessKeysCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.AdminCountAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetAccessKeysCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
