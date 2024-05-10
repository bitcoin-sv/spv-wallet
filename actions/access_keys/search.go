package accesskeys

import (
	"net/http"

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
// @Param		SearchAccessKeys body SearchAccessKeys false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.AccessKey "List of access keys"
// @Failure		400	"Bad request - Error while SearchAccessKeys from request body"
// @Failure 	500	"Internal server error - Error while searching for access keys"
// @Router		/v1/access-key/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams SearchAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessKeys, err := a.Services.SpvWalletEngine.GetAccessKeysByXPubID(
		c.Request.Context(),
		reqXPubID,
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
