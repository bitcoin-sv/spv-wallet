package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// get will get an existing model
// Get access key godoc
// @Summary		Get access key - Use (GET) /api/v1/users/current/keys/{id} instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current/keys/{id} instead.
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id of the access key"
// @Success		200	{object} models.AccessKey "AccessKey with given id"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure		403	"Forbidden - Access key is not owned by the user"
// @Failure 	500	"Internal server error - Error while getting access key"
// @DeprecatedRouter  /v1/access-key [get]
// @Security	x-auth-xpub
func (a *Action) oldGet(c *gin.Context) {
	id := c.Query("id")
	a.getHelper(c, id, true)
}

// get will get an existing model
// Get access key godoc
// @Summary		Get access key
// @Description	Get access key
// @Tags		Access-key
// @Produce		json
// @Param		id path string true "id of the access key"
// @Success		200	{object} response.AccessKey "AccessKey with given id"
// @Failure		400	"Bad request - Missing required field: id"
// @Failure		403	"Forbidden - Access key is not owned by the user"
// @Failure 	500	"Internal server error - Error while getting access key"
// @Router		/api/v1/users/current/keys/{id} [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	id := c.Params.ByName("id")
	a.getHelper(c, id, false)
}

func (a *Action) getHelper(c *gin.Context, id string, snakeCase bool) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	if id == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingFieldID, a.Services.Logger)
		return
	}

	// Get access key
	accessKey, err := a.Services.SpvWalletEngine.GetAccessKey(
		c.Request.Context(), reqXPubID, id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	if accessKey.XpubID != reqXPubID {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, a.Services.Logger)
		return
	}

	if snakeCase {
		contract := mappings.MapToOldAccessKeyContract(accessKey)
		c.JSON(http.StatusOK, contract)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusOK, contract)
}
