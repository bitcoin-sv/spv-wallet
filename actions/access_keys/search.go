package accesskeys

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// search will fetch a list of access keys filtered by metadata
// Search access key godoc
// @Summary		Search access key
// @Description	Search access key
// @Tags		Access-key
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/access-key/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetQueryParameters(c)
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
		fmt.Println("Error getting access keys by xpub id")
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	accessKeyContracts := make([]*models.AccessKey, 0)
	for _, accessKey := range accessKeys {
		accessKeyContracts = append(accessKeyContracts, mappings.MapToAccessKeyContract(accessKey))
	}

	c.JSON(http.StatusOK, accessKeyContracts)
}
