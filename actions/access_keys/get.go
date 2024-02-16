package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
)

// get will get an existing model
// Get access key godoc
// @Summary		Get access key
// @Description	Get access key
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id"
// @Success		200
// @Router		/v1/access-key [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, engine.ErrMissingFieldID)
		return
	}

	// Get access key
	accessKey, err := a.Services.SpvWalletEngine.GetAccessKey(
		c.Request.Context(), reqXPubID, id,
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if accessKey.XpubID != reqXPubID {
		c.JSON(http.StatusForbidden, "unauthorized")
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusOK, contract)
}
