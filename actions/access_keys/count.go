package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
)

// count will fetch a count of access keys filtered by metadata
// Count of access keys godoc
// @Summary		Count of access keys
// @Description	Count of access keys
// @Tags		Access-key
// @Produce		json
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/access-key/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	_, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetAccessKeysByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
