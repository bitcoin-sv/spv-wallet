package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// accessKeysSearch will fetch a list of access keys filtered by metadata
// Access Keys Search godoc
// @Summary		Access Keys Search
// @Description	Access Keys Search
// @Tags		Admin
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @Router		/v1/admin/access-keys/search [post]
// @Security	x-auth-xpub
func (a *Action) accessKeysSearch(c *gin.Context) {
	var reqParams SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessKeys, err := a.Services.SpvWalletEngine.GetAccessKeys(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	accessKeyContracts := make([]*models.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	c.JSON(http.StatusOK, accessKeyContracts)
}

// accessKeysCount will count all access keys filtered by metadata
// Access Keys Count godoc
// @Summary		Access Keys Count
// @Description	Access Keys Count
// @Tags		Admin
// @Produce		json
// @Param		CountRequestParameters body actions.CountRequestParameters false "Enables precise filtering of resource counts using custom conditions or metadata, catering to specific business or analysis needs"
// @Success		200 {number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountRequestParameters from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of access keys"
// @Router		/v1/admin/access-keys/count [post]
// @Security	x-auth-xpub
func (a *Action) accessKeysCount(c *gin.Context) {
	var reqParams CountAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetAccessKeysCount(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
