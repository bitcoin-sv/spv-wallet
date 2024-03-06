package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// search will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key
// @Description	Search access key
// @Tags		Access-key
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "SearchRequestParameters model containing metadata, conditions and query params"
// @Success		200
// @Router		/v1/access-key/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	dbConditions := make(map[string]interface{})
	if conditions != nil {
		dbConditions = *conditions
	}
	dbConditions["xpub_id"] = reqXPubID

	var accessKeys []*engine.AccessKey
	if accessKeys, err = a.Services.SpvWalletEngine.GetAccessKeysByXPubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		&dbConditions,
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
