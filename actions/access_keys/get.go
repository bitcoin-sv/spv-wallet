package accesskeys

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// get will get an existing model
// Get access key godoc
// @Summary		Get access key
// @Description	Get access key
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id of the access key"
// @Success		200	{object} models.AccessKey "AccessKey with given id"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure		403	"Forbidden - Access key is not owned by the user"
// @Failure 	500	"Internal server error - Error while getting access key"
// @Router		/v1/access-key [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	id := c.Query("id")
	if id == "" {
		spverrors2.ErrorResponse(c, spverrors2.ErrMissingFieldID, a.Services.Logger)
		return
	}

	// Get access key
	accessKey, err := a.Services.SpvWalletEngine.GetAccessKey(
		c.Request.Context(), reqXPubID, id,
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	if accessKey.XpubID != reqXPubID {
		spverrors2.ErrorResponse(c, spverrors2.ErrAuthorization, a.Services.Logger)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusOK, contract)
}
