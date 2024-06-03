package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of access keys filtered by metadata
// Count of access keys godoc
// @Summary		Count of access keys
// @Description	Count of access keys
// @Tags		Access-key
// @Produce		json
// @Param		CountAccessKeys body filter.CountAccessKeys false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountAccessKeys from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of access keys"
// @Router		/v1/access-key/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.CountAccessKeys
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetAccessKeysByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
