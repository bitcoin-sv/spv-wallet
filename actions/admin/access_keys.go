package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
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
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "SearchRequestParameters model containing metadata, conditions and query params"
// @Success		200
// @Router		/v1/admin/access-keys/search [post]
// @Security	x-auth-xpub
func (a *Action) accessKeysSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var accessKeys []*engine.AccessKey
	if accessKeys, err = a.Services.SpvWalletEngine.GetAccessKeys(
		c.Request.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
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
// @Param		CountRequestParameters body actions.CountRequestParameters false "CountRequestParameters model containing metadata and conditions"
// @Success		200
// @Router		/v1/admin/access-keys/count [post]
// @Security	x-auth-xpub
func (a *Action) accessKeysCount(c *gin.Context) {
	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetAccessKeysCount(
		c.Request.Context(),
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
